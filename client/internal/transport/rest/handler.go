package rest

import (
	"client/internal/models"
	ws "client/internal/transport/websocket"
	"encoding/json"
	"errors"
	authGrpc "github.com/Smile8MrBread/Chat/auth_service/proto/gen"
	chatGrpc "github.com/Smile8MrBread/Chat/chat_service/proto/gen"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log/slog"
	"math"
	"math/rand/v2"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func StartServer(log *slog.Logger, r *chi.Mux, authClient authGrpc.AuthClient, chatClient chatGrpc.ChatClient) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "client/cmd/frontend/index.html")
	})
	r.Get("/all-chats", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "client/cmd/frontend/all-chats.html")
	})
	r.Get("/contacts", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "client/cmd/frontend/contacts.html")
	})
	//r.Get("/all-groups", func(w http.ResponseWriter, r *http.Request) {
	//	http.ServeFile(w, r, "client/cmd/frontend/all-groups.html")
	//})
	r.Handle("/tmp/*", http.StripPrefix("/tmp/", http.FileServer(http.Dir("client/cmd/tmp"))))

	r.Post("/registration", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			log.Error("Error parsing form", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusConflict)
			jsonResp, _ := json.Marshal(map[string]string{"message": "Internal error"})
			w.Write(jsonResp)
			return
		}

		var localFile string
		file, fileHeader, err := r.FormFile("avatar")
		if err != nil {
			if !errors.Is(err, http.ErrMissingFile) {
				log.Error("Error parsing avatar", slog.String("error", err.Error()))
				w.WriteHeader(http.StatusConflict)
				jsonResp, _ := json.Marshal(map[string]string{"message": "Internal error"})
				w.Write(jsonResp)
				return
			}

			localFile = "default.jpg"
		} else {
			defer file.Close()

			timeNow := strconv.Itoa(int(time.Now().Unix()))
			timeNow += strconv.FormatFloat(math.Pow(rand.Float64(), rand.Float64()), 'E', -1, 32)

			myFile, err := os.Create("client/cmd/tmp/" + timeNow + filepath.Ext(fileHeader.Filename))
			if err != nil {
				log.Error("Error saving img", slog.String("error", err.Error()))
				w.WriteHeader(http.StatusConflict)
				jsonResp, _ := json.Marshal(map[string]string{"message": "Internal error"})
				w.Write(jsonResp)
				return
			}
			defer myFile.Close()

			if _, err = io.Copy(myFile, file); err != nil {
				log.Error("Error copy file", slog.String("error", err.Error()))
				w.WriteHeader(http.StatusConflict)
				jsonResp, _ := json.Marshal(map[string]string{"message": "Internal error"})
				w.Write(jsonResp)
				return
			}

			localFile = timeNow + filepath.Ext(fileHeader.Filename)
		}

		_, err = authClient.Registration(r.Context(), &authGrpc.RegisterRequest{
			Login:     r.FormValue("login"),
			Password:  r.FormValue("password"),
			FirstName: r.FormValue("first-name"),
			LastName:  r.FormValue("last-name"),
			Avatar:    localFile,
		})
		if err != nil {
			if status.Code(err) == codes.AlreadyExists {
				log.Error("User exists", slog.String("error", err.Error()))
				w.WriteHeader(http.StatusConflict)
				jsonResp, _ := json.Marshal(map[string]string{"message": "User exists, login is not unique"})
				w.Write(jsonResp)
				return
			}

			if status.Code(err) == codes.InvalidArgument {
				log.Error("Invalid argument", slog.String("error", err.Error()))
				w.WriteHeader(http.StatusConflict)
				jsonResp, _ := json.Marshal(map[string]string{"message": err.Error()[40:]})
				w.Write(jsonResp)
				return
			}

			log.Error("Error registration", slog.String("login", r.FormValue("login")), slog.String("error", err.Error()))
			w.WriteHeader(http.StatusConflict)
			jsonResp, _ := json.Marshal(map[string]string{"message": "Internal error"})
			w.Write(jsonResp)
			return
		}
	})
	r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		reader := json.NewDecoder(r.Body)
		data := models.RegUser{}

		err := reader.Decode(&data)
		if err != nil {
			log.Error("Error parsing json", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusConflict)
			jsonResp, _ := json.Marshal(map[string]string{"message": "Internal error"})
			w.Write(jsonResp)
			return
		}

		jwtToken, err := authClient.Login(r.Context(), &authGrpc.LoginRequest{
			Login:    data.Login,
			Password: data.Password,
		})
		if err != nil {
			if status.Code(err) == codes.InvalidArgument {
				log.Error("Invalid argument", slog.String("error", err.Error()))
				w.WriteHeader(http.StatusConflict)
				jsonResp, _ := json.Marshal(map[string]string{"message": err.Error()[40:]})
				w.Write(jsonResp)
				return
			}

			log.Error("Error login", slog.String("login", data.Login), slog.String("error", err.Error()))
			w.WriteHeader(http.StatusConflict)
			jsonResp, _ := json.Marshal(map[string]string{"message": "Internal error"})
			w.Write(jsonResp)
			return
		}

		tokenParsed, err := jwt.Parse(jwtToken.Token, func(token *jwt.Token) (interface{}, error) {
			return []byte(models.ReturnSecret()), nil
		})

		claims, ok := tokenParsed.Claims.(jwt.MapClaims)
		if !ok {
			log.Error("Error parsing token", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusConflict)
			jsonResp, _ := json.Marshal(map[string]string{"message": "Internal error"})
			w.Write(jsonResp)
			return
		}
		id := strconv.Itoa(int(claims["user_id"].(float64)))

		http.SetCookie(w, &http.Cookie{
			Name:   "token",
			Value:  jwtToken.Token,
			Path:   "/",
			Secure: false,
		})

		http.SetCookie(w, &http.Cookie{
			Name:   "id",
			Value:  id,
			Path:   "/",
			Secure: false,
		})
	})

	hub := ws.NewHub(log)
	go hub.Run()
	r.HandleFunc("/ws", hub.HandlerWS)
	//r.Handle("/orderbookfeeds", websocket.Handler(server.HandleOrderBook))
	r.Post("/addContact", func(w http.ResponseWriter, r *http.Request) {
		contactid := r.FormValue("conId")

		userid, err := r.Cookie("id")
		if err != nil {
			log.Error("Error getting cookie", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusConflict)
			return
		}

		userId, err := strconv.ParseInt(userid.Value, 10, 64)
		if err != nil {
			log.Error("Error parsing cookie", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusConflict)
			return
		}
		contactId, err := strconv.ParseInt(contactid, 10, 64)
		if err != nil {
			log.Error("Error parsing cookie", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusConflict)
			return
		}

		if _, err = chatClient.AddContact(r.Context(), &chatGrpc.AddContactRequest{
			UserId:    userId,
			ContactId: contactId,
		}); err != nil {
			if status.Code(err) == codes.InvalidArgument {
				log.Error("Invalid argument", slog.String("error", err.Error()))
				w.WriteHeader(http.StatusConflict)
				jsonResp, _ := json.Marshal(map[string]string{"message": err.Error()[40:]})
				w.Write(jsonResp)
				return
			}

			log.Error("Error add contact", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusConflict)
			jsonResp, _ := json.Marshal(map[string]string{"message": "Internal error"})
			w.Write(jsonResp)
			return

		}
	})
	r.Get("/allContacts", func(w http.ResponseWriter, r *http.Request) {
		userid, err := r.Cookie("id")
		if err != nil {
			log.Error("Error getting cookie", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusConflict)
			return
		}

		userId, err := strconv.ParseInt(userid.Value, 10, 64)
		if err != nil {
			log.Error("Error parsing cookie", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusConflict)
			return
		}

		contacts, err := chatClient.AllContacts(r.Context(), &chatGrpc.AllContactsRequest{UserId: userId})
		if err != nil {
			if status.Code(err) == codes.NotFound {
				log.Error("Contacts not found", slog.String("error", err.Error()))
				w.WriteHeader(http.StatusNotAcceptable)
				jsonResp, _ := json.Marshal(map[string]string{"message": err.Error()[40:]})
				w.Write(jsonResp)
				return
			}

			log.Error("Failed to get all contacts", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusConflict)
			return
		}

		buf, err := json.Marshal(contacts)
		if err != nil {
			log.Error("Error marshaling contacts", slog.String("error", err.Error()))
			return
		}

		w.Write(buf)
	})
	r.Get("/isMessaged/{id}", func(w http.ResponseWriter, r *http.Request) {
		strId := chi.URLParam(r, "id")
		toId, err := strconv.ParseInt(strId, 10, 64)
		if err != nil {
			log.Error("Error parse int", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusConflict)
			return
		}

		userId, err := r.Cookie("id")
		if err != nil {
			log.Error("Error getting cookie", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusConflict)
			return
		}

		id, err := strconv.ParseInt(userId.Value, 10, 64)
		if err != nil {
			log.Error("Error parsing cookie", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusConflict)
			return
		}

		_, err = chatClient.IsMessaged(r.Context(), &chatGrpc.IsMessagedRequest{UserId: id, ContactId: toId})
		if err != nil {
			log.Error("Error message", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusConflict)
			return
		}
	})
	r.Get("/allMessaged/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		userId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			log.Error("Error parse id", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusConflict)
			return
		}

		users, err := chatClient.AllMessaged(r.Context(), &chatGrpc.AllMessagedRequest{UserId: userId})
		if err != nil {
			if status.Code(err) == codes.NotFound {
				log.Error("User not found", slog.String("error", err.Error()))
				w.WriteHeader(http.StatusConflict)
				return
			}

			log.Error("Failed to get all messaged", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusConflict)
			return
		}

		data, err := json.Marshal(users)
		if err != nil {
			log.Error("Error encoding json", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusConflict)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	})
	r.Post("/createMessage", func(w http.ResponseWriter, r *http.Request) {
		buf := make([]byte, 1024)
		data := models.Message{}

		n, err := r.Body.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Error("Error reading body", slog.String("error", err.Error()))
				w.WriteHeader(http.StatusConflict)
				return
			}
		}

		err = json.Unmarshal(buf[:n], &data)
		if err != nil {
			log.Error("Error parsing json", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusConflict)
			return
		}

		userId, err := strconv.Atoi(data.UserFrom)
		if err != nil {
			log.Error("Error convert string", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusConflict)
			return
		}
		contactId, err := strconv.Atoi(data.UserTo)
		if err != nil {
			log.Error("Error convert string", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusConflict)
			return
		}

		id, err := chatClient.CreateMessage(r.Context(),
			&chatGrpc.CreateMessageRequest{UserFrom: int64(userId), UserTo: int64(contactId), Date: data.Date, Text: data.Text})
		if err != nil {
			if status.Code(err) == codes.InvalidArgument {
				log.Error("Invalid message", slog.String("error", err.Error()))
				w.WriteHeader(http.StatusConflict)
				return
			}

			log.Error("Failed to create message", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusConflict)
			return
		}

		resp, err := json.Marshal(id)
		if err != nil {
			log.Error("Failed to marshal", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusConflict)
			return
		}

		w.Write(resp)
	})
	r.Get("/identMessage/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		messageId, err := strconv.Atoi(id)
		if err != nil {
			log.Error("Error convert string", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusConflict)
			return
		}

		msg, err := chatClient.IdentMessage(r.Context(), &chatGrpc.IdentMessageRequest{MessageId: int64(messageId)})
		if err != nil {
			if status.Code(err) == codes.InvalidArgument {
				log.Error("Invalid id", slog.String("error", err.Error()))
				w.WriteHeader(http.StatusConflict)
				return
			}

			log.Error("Failed to ident message", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusConflict)
			return
		}

		data, err := json.Marshal(msg)
		if err != nil {
			log.Error("Failed to marshal message", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusConflict)
			return
		}

		w.Write(data)
	})
	r.Get("/allMessages/{userFrom}/{userTo}", func(w http.ResponseWriter, r *http.Request) {
		userFrom := chi.URLParam(r, "userFrom")
		userTo := chi.URLParam(r, "userTo")

		usrFrom, err := strconv.Atoi(userFrom)
		if err != nil {
			log.Error("Error convert string", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusConflict)
			return
		}
		usrTo, err := strconv.Atoi(userTo)
		if err != nil {
			log.Error("Error convert string", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusConflict)
			return
		}

		ids, err := chatClient.AllMessages(r.Context(), &chatGrpc.AllMessagesRequest{
			UserFrom: int64(usrFrom),
			UserTo:   int64(usrTo),
		})
		if err != nil {
			if status.Code(err) == codes.NotFound {
				log.Error("User not found", slog.String("error", err.Error()))
				w.WriteHeader(http.StatusConflict)
				return
			}

			log.Error("Internal error", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusConflict)
			return
		}

		data, err := json.Marshal(ids.GetMessageIds())
		if err != nil {
			log.Error("Error marshaling", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusConflict)
			return
		}

		w.Write(data)
	})

	r.Get("/identity/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		userId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			log.Error("Error ident user", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusConflict)
			return
		}
		user, err := authClient.IdentUser(r.Context(), &authGrpc.IdentUserRequest{UserId: userId})
		if err != nil {
			if status.Code(err) == codes.NotFound {
				log.Error("Error user not found", slog.String("error", err.Error()))
				w.WriteHeader(http.StatusConflict)
				jsonResp, _ := json.Marshal(map[string]string{"message": err.Error()[40:]})
				w.Write(jsonResp)
				return
			}

			log.Error("Error ident user", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusConflict)
			return
		}
		b, err := json.Marshal(user)
		if err != nil {
			if err != nil {
				log.Error("Error encoding json", slog.String("error", err.Error()))
				w.WriteHeader(http.StatusConflict)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	})

	http.ListenAndServe(":8080", r)
}
