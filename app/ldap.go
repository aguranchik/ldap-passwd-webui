package app

import (
//      "crypto/tls"
//      "fmt"
//      "log"
//      "strings"
//
//      "gopkg.in/ldap.v2"
        auth "github.com/korylprince/go-ad-auth/v3"
        "log"
        "fmt"
//        "os"
        ss "strings"
)

func bindUser(l *auth.Conn, userDN, passwd string) (bool, error) {
        log.Printf("\nBinding with userDN: %s", userDN)
        status, err := l.Bind(userDN, passwd)
        if err != nil {
                return status, fmt.Errorf("Dial: %v", err)
        }
        log.Printf("\nBound successfully with userDN: %s", userDN)
        return status, nil
}

//func bindUser(l *ldap.Conn, userDN, passwd string) error {
//      log.Printf("\nBinding with userDN: %s", userDN)
//      err := l.Bind(userDN, passwd)
//      if err != nil {
//              log.Printf("\nLDAP auth. failed for %s, reason: %v", userDN, err)
//              return err
//      }
//      log.Printf("\nBound successfully with userDN: %s", userDN)
//      return err
//}

func sanitizedUserDN(username string) (string, bool) {
        badCharacters := "\x00()*\\,='\"#+;<>"
        userDN := ""
        if ss.ContainsAny(username, badCharacters) {
                log.Printf("\n'%s' contains invalid DN characters. Aborting.", username)
                return userDN, false
        }

        before, after, found := ss.Cut(username, "@")
        if found {
                if after == "21vek.by" {
                        userDN := before + "@21vek.local"
                        log.Printf("\nUsername changed to: %s", userDN)
                        return userDN, true
                } else if after == "21vek.local" {
                        userDN := before + "@21vek.local"
                        log.Printf("\nUsername entered correctly: %s", userDN)
                        return userDN, true
                } else {
                        log.Printf("\nInvalid domain name: %s",after)
                        return userDN, false
                }
        } else {
                userDN := username + "@21vek.local"
                log.Printf("\nUsername changed to: ", userDN)
                return userDN, true
        }
        return userDN, true
}

//func (ls *LDAPClient) sanitizedUserDN(username string) (string, bool) {
        // See http://tools.ietf.org/search/rfc4514: "special characters"
//      badCharacters := "\x00()*\\,='\"#+;<>"
//      if strings.ContainsAny(username, badCharacters) {
//              log.Printf("\n'%s' contains invalid DN characters. Aborting.", username)
//              return "", false
//      }
//
//      return fmt.Sprintf(ls.UserDN, username), true
//}

func dial(ls *auth.Config) (*auth.Conn, error) {
        log.Printf("\nDialing LDAP with security protocol: ", ls.Security)
        conn, err := ls.Connect()
        if err != nil {
                return nil, fmt.Errorf("Dial: %v", err)
        }
        return conn, nil
}


//func dial(ls *LDAPClient) (*ldap.Conn, error) {
//      log.Printf("\nDialing LDAP with security protocol (%v) without verifying: %v", ls.SecurityProtocol, ls.SkipVerify)
//
//      tlsCfg := &tls.Config{
//              ServerName:         ls.Host,
//              InsecureSkipVerify: ls.SkipVerify,
//      }
//      if ls.SecurityProtocol == SecurityProtocolLDAPS {
//              return ldap.DialTLS("tcp", fmt.Sprintf("%s:%d", ls.Host, ls.Port), tlsCfg)
//      }
//
//      conn, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", ls.Host, ls.Port))
//      if err != nil {
//              return nil, fmt.Errorf("Dial: %v", err)
//      }
//
//      if ls.SecurityProtocol == SecurityProtocolStartTLS {
//              if err = conn.StartTLS(tlsCfg); err != nil {
//                      conn.Close()
//                      return nil, fmt.Errorf("StartTLS: %v", err)
//              }
//      }
//
//      return conn, nil
//}

func ModifyPassword(ls *auth.Config, name, passwd, newPassword string) error {
        if len(passwd) == 0 {
                return fmt.Errorf("Auth failed for %s, password cannot be empty", name)
        }
        l, err := dial(ls)
        if err != nil {
                return fmt.Errorf("LDAP Connect error, %s:%v", ls.Server, err)
        } else {
                log.Printf("\nСonnection successful!")
        }

        var userDN string
        log.Printf("\nLDAP will bind directly via: %s", name)

        var ok bool
        userDN, ok = sanitizedUserDN(name)
        if !ok {
                return fmt.Errorf("Error sanitizing name %s", name)
        }
        _, err2 := bindUser(l, userDN, passwd)
        if err2 != nil {
                return fmt.Errorf("Auth for %s in server %s error: %v", userDN, ls.Server, err2)
        }

        log.Printf("\nLDAP will execute password change on: %s", userDN)
        err3 := auth.UpdatePassword(ls, userDN, passwd, newPassword)
        if err3 != nil {
                log.Printf("\nError execute password change on: %s", err3)
        }

        return err3
}


// ModifyPassword : modify user's password
//func (ls *LDAPClient) ModifyPassword(name, passwd, newPassword string) error {
//      if len(passwd) == 0 {
//              return fmt.Errorf("Auth. failed for %s, password cannot be empty", name)
//      }
//      l, err := dial(ls)
//      if err != nil {
//              ls.Enabled = false
//              return fmt.Errorf("LDAP Connect error, %s:%v", ls.Host, err)
//      }
//      defer l.Close()
//
//      var userDN string
//      log.Printf("\nLDAP will bind directly via UserDN template: %s", ls.UserDN)
//
//      var ok bool
//      userDN, ok = ls.sanitizedUserDN(name)
//      if !ok {
//              return fmt.Errorf("Error sanitizing name %s", name)
//      }
//      bindUser(l, userDN, passwd)
//
//      log.Printf("\nLDAP will execute password change on: %s", userDN)
//      req := ldap.NewPasswordModifyRequest(userDN, passwd, newPassword)
//      _, err = l.PasswordModify(req)
//
//      return err
//}

func NewLDAPServer() *auth.Config {
        return &auth.Config{
                Server:           envStr("LPW_HOST", "01dc21vek.21vek.local"),
                Port:             envInt("LPW_Port", 3899),
                Security:         auth.SecurityStartTLS,
                BaseDN:           envStr("LPW_BaseDN", "DC=21vek,DC=local"),
        }
}



// NewLDAPClient : Creates new LDAPClient capable of binding and changing passwords
//func NewLDAPClient() *LDAPClient {
//
//      securityProtocol := SecurityProtocolUnencrypted
//      if envBool("LPW_ENCRYPTED", true) {
//              securityProtocol = SecurityProtocolLDAPS
//              if envBool("LPW_START_TLS", false) {
//                      securityProtocol = SecurityProtocolStartTLS
//              }
//      }
//
//      return &LDAPClient{
//              Host:             envStr("LPW_HOST", ""),
//              Port:             envInt("LPW_PORT", 636), // 389
//              SecurityProtocol: securityProtocol,
//              SkipVerify:       envBool("LPW_SSL_SKIP_VERIFY", false),
//              UserDN:           envStr("LPW_USER_DN", "uid=%s,ou=people,dc=example,dc=org"),
//              UserBase:         envStr("LPW_USER_BASE", "ou=people,dc=example,dc=org"),
//      }
//}
