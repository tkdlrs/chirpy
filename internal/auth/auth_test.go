package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCheckPasswordHash(t *testing.T) {
	//
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name     string
		password string
		hash     string
		wantErr  bool
	}{
		{
			name:     "Correct password",
			password: password1,
			hash:     hash1,
			wantErr:  false,
		},
		{
			name:     "Incorrect password",
			password: "wrongPassword",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Password doesn't match different hash",
			password: password1,
			hash:     hash2,
			wantErr:  true,
		},
		{
			name:     "Empty password",
			password: "",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Invalid hash",
			password: password1,
			hash:     "invalidhash",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Messing with testing (since no one want to take the time to explain it to me...)
func TestCheckTrueIsTrue(t *testing.T) {
	t.Run("Sanity Check", func(t *testing.T) {
		if 2+2 != 4 {
			t.Error("two plus two should be equal to four")
		}
	})
}

// testing the JWT stuff
func TestCheckJWT(t *testing.T) {
	//
	userID := uuid.MustParse("cefe4238-72ea-4aa5-b154-d7048777e572")
	tests := []struct {
		name              string
		userid            uuid.UUID
		secret            string
		time              time.Duration
		wantErrMaking     bool
		wantErrValidating bool
	}{
		{
			name:              "JWT check -happy path",
			userid:            userID,
			secret:            "secret",
			time:              1 * time.Hour,
			wantErrMaking:     false,
			wantErrValidating: false,
		},
		// {
		// 	name:              "JWT check -expired token",
		// 	userid:            userID,
		// 	secret:            "Secret",
		// 	time:              1 * time.Millisecond,
		// 	wantErrMaking:     false,
		// 	wantErrValidating: true,
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jwtToken, err := MakeJWT(tt.userid, tt.secret, tt.time)
			if (err != nil) != tt.wantErrMaking {
				t.Errorf("MakeJWT() error = %v, wantErrMaking %v", err, tt.wantErrMaking)
			}
			//
			time.Sleep(1 * time.Second)
			//
			_, err = ValidateJWT(jwtToken, tt.secret)
			if (err != nil) != tt.wantErrValidating {
				t.Errorf("ValidateJWT() error = %v, wantErrValidating %v", err, tt.wantErrValidating)
			}
			//
			// if userUUID != userID {
			// 	t.Error("The userUUID that we got back from validation does not equal the original userID.")
			// }

		})
	}
}
