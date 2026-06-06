package domain

import "time"

type User struct {
	ID           int
	Username     string
	Email        string
	Password 		 string
	RegisterAt   time.Time
}

type Shortlink struct {
	ID          int
	UserID      int
	Shortlink   string
	OriginalUrl string
	Count       int
	CreatedAt   time.Time
}

type ShortlinkInfo struct {
	ID           int
	IPAddr       string
	UserAgent    string
	ShortlinkID  int
	LastAccessAt time.Time
}

type JWT struct {
	UserID int
	Token  string
}
type Session struct {
	UserID       int
	RefreshToken string
	Device       *Device
}

type Device struct {
	IpAddr    string
	UserAgent string
}

func NewShortLink(
	userID int,
	shortLink string,
	originalLink string,
) *Shortlink {
	return &Shortlink{
		UserID:      userID,
		Shortlink:   shortLink,
		OriginalUrl: originalLink,
	}
}

func NewJwt(
	userID int,
	token string,
) *JWT {
	return &JWT{
		UserID: userID,
		Token:  token,
	}
}

func NewUser(
	username string,
	email string,
	password string,
) *User {
	return &User{
		Username:     username,
		Email:        email,
		Password: password,
	}
}

func NewLoginUser(
	username string,
	password string,
) *User {
	return &User{
		Username:     username,
		Password: password,
	}
}

func NewDevice(
	ipAddr string,
	userAgent string,
) *Device {
	return &Device{
		IpAddr:    ipAddr,
		UserAgent: userAgent,
	}
}

func NewSession(
	userID int,
	refreshToken string,
	device *Device,
) *Session {
	return &Session{
		UserID:       userID,
		RefreshToken: refreshToken,
		Device:       device,
	}
}

func NewShortLinkInfo(
	ip_addr string,
	userAgent string,
	shortlinkID int,
) *ShortlinkInfo {
	return &ShortlinkInfo{
		IPAddr:      ip_addr,
		UserAgent:   userAgent,
		ShortlinkID: shortlinkID,
	}
}
