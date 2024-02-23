package verifier

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

func TestOptions(t *testing.T) {
	tests := []struct {
		name string
		opts []Options
		want *options
	}{{
		name: "no options",
		want: &options{},
	}, {
		name: "signature option",
		opts: []Options{WithPublicKey([]byte("foo"))},
		want: &options{
			PublicKey: []byte("foo"),
			ROpt:      nil,
		},
	}, {
		name: "keychain option",
		opts: []Options{WithRemoteOptions(remote.WithAuthFromKeychain(authn.DefaultKeychain))},
		want: &options{
			PublicKey: nil,
			ROpt:      []remote.Option{remote.WithAuthFromKeychain(authn.DefaultKeychain)},
		},
	}, {
		name: "keychain and authenticator option",
		opts: []Options{WithRemoteOptions(
			remote.WithAuth(&authn.Basic{Username: "foo", Password: "bar"}),
			remote.WithAuthFromKeychain(authn.DefaultKeychain),
		)},
		want: &options{
			PublicKey: nil,
			ROpt: []remote.Option{
				remote.WithAuth(&authn.Basic{Username: "foo", Password: "bar"}),
				remote.WithAuthFromKeychain(authn.DefaultKeychain),
			},
		},
	}, {
		name: "keychain, authenticator and transport option",
		opts: []Options{WithRemoteOptions(
			remote.WithAuth(&authn.Basic{Username: "foo", Password: "bar"}),
			remote.WithAuthFromKeychain(authn.DefaultKeychain),
			remote.WithTransport(http.DefaultTransport),
		)},
		want: &options{
			PublicKey: nil,
			ROpt: []remote.Option{
				remote.WithAuth(&authn.Basic{Username: "foo", Password: "bar"}),
				remote.WithAuthFromKeychain(authn.DefaultKeychain),
				remote.WithTransport(http.DefaultTransport),
			},
		},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := options{}
			for _, opt := range test.opts {
				opt(&o)
			}
			if !reflect.DeepEqual(o.PublicKey, test.want.PublicKey) {
				t.Errorf("got %#v, want %#v", &o.PublicKey, test.want.PublicKey)
			}

			if test.want.ROpt != nil {
				if len(o.ROpt) != len(test.want.ROpt) {
					t.Errorf("got %d remote options, want %d", len(o.ROpt), len(test.want.ROpt))
				}
				return
			}

			if test.want.ROpt == nil {
				if len(o.ROpt) != 0 {
					t.Errorf("got %d remote options, want %d", len(o.ROpt), 0)
				}
			}
		})
	}
}
