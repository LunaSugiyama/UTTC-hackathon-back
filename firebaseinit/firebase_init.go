package firebaseinit

import (
	"context"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

var AuthClient *auth.Client

func InitFirebase() error {
	opt := option.WithCredentialsFile("/home/denjo/ダウンロード/term4-luna-sugiyama-firebase-adminsdk-1joai-b0f371c4d8.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	// app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		return err
	}

	AuthClient, err = app.Auth(context.Background())
	if err != nil {
		return err
	}

	return nil
}

type FirebaseConfig struct {
	Type                string `json:"type"`
	ProjectID           string `json:"project_id"`
	PrivateKeyID        string `json:"private_key_id"`
	PrivateKey          string `json:"private_key"`
	ClientEmail         string `json:"client_email"`
	ClientID            string `json:"client_id"`
	AuthURI             string `json:"auth_uri"`
	TokenURI            string `json:"token_uri"`
	AuthProviderCertURL string `json:"auth_provider_x509_cert_url"`
	ClientCertURL       string `json:"client_x509_cert_url"`
	UniverseDomain      string `json:"universe_domain"`
}

// func InitFirebase() error {
// err := godotenv.Load(".env")
// if err != nil {
// 	log.Fatal("Error loading .env file")
// 	return err
// }

// config := FirebaseConfig{
// 	Type:                os.Getenv("TYPE"),
// 	ProjectID:           os.Getenv("PROJECT_ID"),
// 	PrivateKeyID:        os.Getenv("PRIVATE_KEY_ID"),
// 	PrivateKey:          os.Getenv("PRIVATE_KEY"),
// 	ClientEmail:         os.Getenv("CLIENT_EMAIL"),
// 	ClientID:            os.Getenv("CLIENT_ID"),
// 	AuthURI:             os.Getenv("AUTH_URI"),
// 	TokenURI:            os.Getenv("TOKEN_URI"),
// 	AuthProviderCertURL: os.Getenv("AUTH_PROVIDER_X509_CERT_URL"),
// 	ClientCertURL:       os.Getenv("CLIENT_X509_CERT_URL"),
// 	UniverseDomain:      os.Getenv("UNIVERSE_DOMAIN"),
// }

// log.Printf("config: %v", config)

// configBytes, err := json.Marshal(config)
// if err != nil {
// 	log.Fatal("Error marshaling FirebaseConfig to JSON")
// 	return err
// }

// log.Printf("configBytes: %v", configBytes)

// opt := option.WithCredentialsJSON(configBytes)
// 	app, err := firebase.NewApp(context.Background(), nil)
// 	if err != nil {
// 		log.Fatalf("error initializing app: %v\n", err)
// 	}

// 	AuthClient, err = app.Auth(context.Background())
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
