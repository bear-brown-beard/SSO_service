package models

import (
	"database/sql"
)

type User struct {
	ID                 int64          `db:"id" json:"id"`
	Username           string         `db:"username" json:"username"`
	AuthKey            string         `db:"auth_key" json:"auth_key"`
	Lang               string         `db:"lang" json:"lang"`
	DeviceToken        sql.NullString `db:"device_token" json:"device_token,omitempty"`
	PasswordHash       string         `db:"password_hash" json:"-"`
	PasswordResetToken sql.NullString `db:"password_reset_token" json:"-"`
	Email              sql.NullString `db:"email" json:"email,omitempty"`
	Status             int16          `db:"status" json:"status"`
	CreatedAt          int64          `db:"created_at" json:"created_at"`
	UpdatedAt          int64          `db:"updated_at" json:"updated_at"`
	LoginAt            sql.NullInt64  `db:"login_at" json:"login_at,omitempty"`
	FirstName          sql.NullString `db:"first_name" json:"first_name,omitempty"`
	LastName           sql.NullString `db:"last_name" json:"last_name,omitempty"`
	Address            sql.NullString `db:"address" json:"address,omitempty"`
	Phone              sql.NullString `db:"phone" json:"phone,omitempty"`
	CityID             sql.NullInt64  `db:"city_id" json:"city_id,omitempty"`
	Avatar             sql.NullString `db:"avatar" json:"avatar,omitempty"`
	IP                 sql.NullString `db:"ip" json:"ip,omitempty"`
	Browser            sql.NullString `db:"browser" json:"browser,omitempty"`
	IsGuest            sql.NullBool   `db:"is_guest" json:"is_guest,omitempty"`
	PlayerID           sql.NullString `db:"player_id" json:"player_id,omitempty"`
	SignalCarts        sql.NullBool   `db:"signal_carts" json:"signal_carts,omitempty"`
	SignalOrders       sql.NullBool   `db:"signal_orders" json:"signal_orders,omitempty"`
	PharmacyID         sql.NullInt64  `db:"pharmacy_id" json:"pharmacy_id,omitempty"`
	DeviceID           sql.NullString `db:"device_id" json:"device_id,omitempty"`
	DeletedAt          sql.NullTime   `db:"deleted_at" json:"deleted_at,omitempty"`
	LastVisitedAt      sql.NullTime   `db:"last_visited_at" json:"last_visited_at,omitempty"`
}

type TestAccount struct {
	ID    int64          `db:"id" json:"id"`
	Name  sql.NullString `db:"name" json:"name,omitempty"`
	Phone string         `db:"phone" json:"phone"`
	Code  sql.NullString `db:"code" json:"code,omitempty"`
}

type UserAccessToken struct {
	ID        int64          `db:"id" json:"id"`
	UserID    sql.NullInt64  `db:"user_id" json:"user_id,omitempty"`
	Token     sql.NullString `db:"token" json:"token,omitempty"`
	Agent     sql.NullString `db:"agent" json:"agent,omitempty"`
	IP        sql.NullString `db:"ip" json:"ip,omitempty"`
	ExpireAt  sql.NullTime   `db:"expire_at" json:"expire_at,omitempty"`
	LoginAt   sql.NullTime   `db:"login_at" json:"login_at,omitempty"`
	CreatedAt sql.NullTime   `db:"created_at" json:"created_at,omitempty"`
	UpdatedAt sql.NullTime   `db:"updated_at" json:"updated_at,omitempty"`
}

type UserMindBox struct {
	ID                     int64          `db:"id" json:"id"`
	UserID                 sql.NullInt64  `db:"user_id" json:"user_id,omitempty"`
	ConsentToMailings      sql.NullBool   `db:"consent_to_mailings" json:"consent_to_mailings,omitempty"`
	LoyaltyProgramEnrolled sql.NullBool   `db:"loyalty_program_enrolled" json:"loyalty_program_enrolled,omitempty"`
	MindBoxUserID          sql.NullString `db:"mind_box_user_id" json:"mind_box_user_id,omitempty"`
	IsPhoneConfirm         sql.NullBool   `db:"is_phone_confirm" json:"is_phone_confirm,omitempty"`
	Barcode                sql.NullString `db:"barcode" json:"barcode,omitempty"`
	RefPromo               sql.NullString `db:"ref_promo" json:"ref_promo,omitempty"`
	RefPromoPharm          sql.NullString `db:"ref_promo_pharm" json:"ref_promo_pharm,omitempty"`
}

type SMSCResponse struct {
	ID     int    `json:"id"`
	Count  int    `json:"count"`
	Cost   string `json:"cost"`
	Status int    `json:"status"`
	Error  string `json:"error"`
}
