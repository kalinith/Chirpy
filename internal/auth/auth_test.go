package auth
import "testing"

func TestHashing(t *testing.T) {
	cases := []struct {
		input	string
	}{
		{
			input:    "hello world",
		},
                {
                        input:    "This is a test",
                },
                {
                        input:    "the Quick brown Fox",
                },
                {
                        input:    "jumped over the lazy dog",
                },
                {
                        input:    "Uno Duo Tres Quatro Sinco sinco ses",
                },
                {
                        input:    "thisisareallylongsentancewithnospaces",
                },
                {
                        input:    "oh",
                },
		// add more cases here
	}

	for _, c := range cases {
		hash, err := HashPassword(c.input)
		// get a hash for the password
		if err != nil {
			t.Errorf("failed to hash password '%s': %s",c.input, err)
		}
		// test if password can be compared to hash
		passmatch := CheckPasswordHash(hash, c.input)
		if passmatch != nil {
			t.Errorf("password(%s) doesn't match Hash(%s)", c.input, hash)
		}
	}

}
