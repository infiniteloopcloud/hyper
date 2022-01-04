package hyper

import (
	"encoding/base64"
	"testing"
)

func TestGetKeyPair(t *testing.T) {
	expectedCert := "MIIBhTCCASugAwIBAgIQIRi6zePL6mKjOipn+dNuaTAKBggqhkjOPQQDAjASMRAwDgYDVQQKEwdBY21lIENvMB4XDTE3MTAyM" +
		"DE5NDMwNloXDTE4MTAyMDE5NDMwNlowEjEQMA4GA1UEChMHQWNtZSBDbzBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABD0d7VNhb" +
		"WvZLWPuj/RtHFjvtJBEwOkhbN/BnnE8rnZR8+sbwnc/KhCk3FhnpHZnQz7B5aETbbIgmuvewdjvSBSjYzBhMA4GA1UdDwEB/wQEA" +
		"wICpDATBgNVHSUEDDAKBggrBgEFBQcDATAPBgNVHRMBAf8EBTADAQH/MCkGA1UdEQQiMCCCDmxvY2FsaG9zdDo1NDUzgg4xMjcuM" +
		"C4wLjE6NTQ1MzAKBggqhkjOPQQDAgNIADBFAiEA2zpJEPQyz6/lWf86aX6PepsntZv2GYlA5UpabfT2EZICICpJ5h/iI+i341gBm" +
		"LiAFQOyTDT+/wQc6MF9+Yw1Yy0t"

	var certManager = environmentTLS{
		EnvironmentTLSOpts{
			TLSCert:          "t_cert",
			TLSCertBlockType: "t_cert_block",
			TLSKey:           "t_key",
			TLSKeyBlockType:  "t_key_block",
		},
	}

	getenv = func(key string) string {
		switch key {
		case certManager.env.TLSCert:
			return "MIIBhTCCASugAwIBAgIQIRi6zePL6mKjOipn+dNuaTAKBggqhkjOPQQDAjASMRAwDgYDVQQKEwdBY21lIENvMB4XDTE3MTAyM" +
				"DE5NDMwNloXDTE4MTAyMDE5NDMwNlowEjEQMA4GA1UEChMHQWNtZSBDbzBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABD0d7VNhb" +
				"WvZLWPuj/RtHFjvtJBEwOkhbN/BnnE8rnZR8+sbwnc/KhCk3FhnpHZnQz7B5aETbbIgmuvewdjvSBSjYzBhMA4GA1UdDwEB/wQEA" +
				"wICpDATBgNVHSUEDDAKBggrBgEFBQcDATAPBgNVHRMBAf8EBTADAQH/MCkGA1UdEQQiMCCCDmxvY2FsaG9zdDo1NDUzgg4xMjcuM" +
				"C4wLjE6NTQ1MzAKBggqhkjOPQQDAgNIADBFAiEA2zpJEPQyz6/lWf86aX6PepsntZv2GYlA5UpabfT2EZICICpJ5h/iI+i341gBm" +
				"LiAFQOyTDT+/wQc6MF9+Yw1Yy0t"
		case certManager.env.TLSCertBlockType:
			return "CERTIFICATE"
		case certManager.env.TLSKey:
			return "MHcCAQEEIIrYSSNQFaA2Hwf1duRSxKtLYX5CB04fSeQ6tF1aY/PuoAoGCCqGSM49AwEHoUQDQgAEPR3tU2Fta9ktY+6P9G0cW" +
				"O+0kETA6SFs38GecTyudlHz6xvCdz8qEKTcWGekdmdDPsHloRNtsiCa697B2O9IFA=="
		case certManager.env.TLSKeyBlockType:
			return "EC PRIVATE KEY"
		default:
			return ""
		}
	}

	cfg, err := certManager.from()
	if err != nil {
		t.Fatal(err)
	}

	if cert := base64.StdEncoding.EncodeToString(cfg.Certificates[0].Certificate[0]); cert != expectedCert {
		t.Errorf("Expected cert should be:\n %s\n instead of:\n%s\n", expectedCert, cert)
	}
}
