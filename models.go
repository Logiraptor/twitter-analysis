package main

type Status struct {
	CreatedAt            string   `json:"created_at"`
	ID                   int64    `json:"id"`
	IDStr                string   `json:"id_str"`
	Text                 string   `json:"text"`
	Source               string   `json:"source"`
	Truncated            bool     `json:"truncated"`
	InReplyToStatusID    int64    `json:"in_reply_to_status_id"`
	InReplyToStatusIDStr string   `json:"in_reply_to_status_id_str"`
	InReplyToUserID      int64    `json:"in_reply_to_user_id"`
	InReplyToUserIDStr   string   `json:"in_reply_to_user_id_str"`
	InReplyToScreenName  string   `json:"in_reply_to_screen_name"`
	User                 user     `json:"user"`
	Geo                  *geo     `json:"geo"`
	Coordinates          *geo     `json:"coordinates"`
	Place                *place   `json:"place"`
	Contributors         []int64  `json:"contributors"`
	RetweetCount         int      `json:"retweet_count"`
	FavoriteCount        int      `json:"favorite_count"`
	Entities             entities `json:"entities"`
	Favorited            bool     `json:"favorited"`
	Retweeted            bool     `json:"retweeted"`
	FilterLevel          string   `json:"filter_level"`
	Lang                 string   `json:"lang"`
	ComputedCoords       *geo     `json:"computed_coords"`
}

type place struct {
	Attributes struct {
	} `json:"attributes"`
	ID          string `json:"id"`
	URL         string `json:"url"`
	PlaceType   string `json:"place_type"`
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	CountryCode string `json:"country_code"`
	Country     string `json:"country"`
	BoundingBox struct {
		Coordinates [][][]float64 `json:"coordinates"`
		Type        string        `json:"type"`
	} `json:"bounding_box"`
}

type entities struct {
	URLs         []url     `json:"urls"`
	Hashtags     []hashtag `json:"hashtags"`
	UserMentions []user    `json:"user_mentions"`
}
type geo struct {
	Type        string     `json:"type"`
	Coordinates [2]float32 `json:"coordinates"`
}
type hashtag struct {
	Text    string `json:"text"`
	Indices []int  `json:"indices"`
}
type url struct {
	URL         string `json:"url"`
	DisplayURL  string `json:"display_url"`
	ExpandedURL string `json:"expanded_url"`
	Indices     []int  `json:"indices"`
}
type user struct {
	ID                             int64   `json:"id"`
	IDStr                          string  `json:"id_str"`
	Name                           string  `json:"name"`
	ScreenName                     string  `json:"screen_name"`
	Location                       *string `json:"location"`
	URL                            string  `json:"url"`
	Description                    string  `json:"string"`
	Protected                      bool    `json:"protected"`
	FollowersCount                 int     `json:"followers_count"`
	FriendsCount                   int     `json:"friends_count"`
	ListedCount                    int     `json:"listed_count"`
	CreatedAt                      string  `json:"created_at"`
	FavouritesCount                int     `json:"favourites_count"`
	UtcOffset                      int     `json:"utc_offset"`
	TimeZone                       string  `json:"time_zone"`
	GeoEnabled                     bool    `json:"geo_enabled"`
	Verified                       bool    `json:"verified"`
	StatusesCount                  int64   `json:"statuses_count"`
	Lang                           string  `json:"lang"`
	ContributorsEnabled            bool    `json:"contributors_enabled"`
	IsTranslator                   bool    `json:"is_translator"`
	IsTranslationEnabled           bool    `json:"is_translation_enabled"`
	ProfileBackgroundColor         string  `json:"profile_background_color"`
	ProfileBackgroundImageURL      string  `json:"profile_background_image_url"`
	ProfileBackgroundImageURLHttps string  `json:"profile_background_image_url_https"`
	ProfileBackgroundTile          bool    `json:"profile_background_tile"`
	ProfileImageURL                string  `json:"profile_image_url"`
	ProfileImageURLHttps           string  `json:"profile_image_url_https"`
	ProfileLinkColor               string  `json:"profile_link_color"`
	ProfileSidebarBorderColor      string  `json:"profile_sidebar_border_color"`
	ProfileSidebarFillColor        string  `json:"profile_sidebar_fill_color"`
	ProfileTextColor               string  `json:"profile_text_color"`
	ProfileUseBackgroundImage      bool    `json:"profile_use_background_image"`
	DefaultProfile                 bool    `json:"default_profile"`
	DefaultProfileImage            bool    `json:"default_profile_image"`
	Following                      bool    `json:"following"`
	FollowRequestSent              bool    `json:"follow_request_sent"`
	Notifications                  bool    `json:"notifications"`
}
