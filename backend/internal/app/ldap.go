package app

import (
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"gopkg.in/ldap.v3"
)

type User struct {
	Cid    string   `json:"cid"`
	Nick   string   `json:"nick"`
	Groups []string `json:"groups"`
	jwt.StandardClaims
}

func login_ldap(cid string, pass string) (User, error) {
	conn, err := ldap.DialURL("ldap://kamino.chalmers.it")
	if err != nil {
		return User{}, errors.New("Could not connect to ldap")
	}
	defer conn.Close()

	err = conn.Bind(fmt.Sprintf("uid=%s,ou=people,dc=chalmers,dc=it", cid), pass)
	if err != nil {
		return User{}, errors.New("Username and password did not match")
	}

	user := User{Cid: cid, Groups: []string{}}

	s, er := search(conn,
		"ou=people,dc=chalmers,dc=it",
		fmt.Sprintf("(uid=%s)", cid),
		[]string{"nickname"})

	if er == nil {
		user.Nick = s.Entries[0].GetAttributeValue("nickname")
	}

	user.Groups, _ = getGroups(conn,
		fmt.Sprintf("uid=%s,ou=people,dc=chalmers,dc=it", cid))

	return user, nil
}

func search(conn *ldap.Conn, dc string, filter string, attributes []string) (*ldap.SearchResult, error) {
	return conn.Search(ldap.NewSearchRequest(
		dc,
		ldap.ScopeWholeSubtree, ldap.DerefAlways, 0, 0, false,
		filter,
		attributes,
		nil))
}

func getGroups(conn *ldap.Conn, member string) ([]string, error) {
	membersQuery := func(dn string) (*ldap.SearchResult, error) {
		return search(conn, "ou=groups,dc=chalmers,dc=it",
			fmt.Sprintf("(&(|(objectClass=itGroup)(objectClass=itPosition))(member=%s))", dn),
			[]string{"cn"})
	}

	s, err := membersQuery(member)
	if err != nil {
		return nil, err
	}

	groups := []string{}

	for _, entry := range s.Entries {
		groups = append(groups, entry.GetAttributeValue("cn"))

		//If the group is member of the root group
		if s, _ := membersQuery(entry.DN); len(s.Entries) != 0 {
			groups = append(groups, s.Entries[0].GetAttributeValue("cn"))
		}
	}
	return groups, nil
}
