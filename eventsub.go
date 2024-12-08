package helix

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// EventSub Types for Parsing Requests / Responses

// Represents a subscription
type EventSubSubscription struct {
	ID        string            `json:"id"`
	Type      string            `json:"type"`
	Version   string            `json:"version"`
	Status    string            `json:"status"`
	Condition EventSubCondition `json:"condition"`
	Transport EventSubTransport `json:"transport"`
	CreatedAt Time              `json:"created_at"`
	Cost      int               `json:"cost"`
}

// Conditions for a subscription, not all are necessary and some only apply to some subscription types, see https://dev.twitch.tv/docs/eventsub/eventsub-reference
type EventSubCondition struct {
	BroadcasterUserID     string `json:"broadcaster_user_id"`
	FromBroadcasterUserID string `json:"from_broadcaster_user_id"`
	ModeratorUserID       string `json:"moderator_user_id"`
	ToBroadcasterUserID   string `json:"to_broadcaster_user_id"`
	RewardID              string `json:"reward_id"`
	ClientID              string `json:"client_id"`
	ExtensionClientID     string `json:"extension_client_id"`
	UserID                string `json:"user_id"`
}

// Transport for the subscription, currently the only supported Method is "webhook". Secret must be between 10 and 100 characters
type EventSubTransport struct {
	Method    string `json:"method"`
	Callback  string `json:"callback"`
	Secret    string `json:"secret"`
	SessionID string `json:"session_id"`
}

// Twitch Response for getting all current subscriptions
type ManyEventSubSubscriptions struct {
	Total                 int                    `json:"total"`
	TotalCost             int                    `json:"total_cost"`
	MaxTotalCost          int                    `json:"max_total_cost"`
	EventSubSubscriptions []EventSubSubscription `json:"data"`
	Pagination            Pagination             `json:"pagination"`
}

// Response for getting all current subscriptions
type EventSubSubscriptionsResponse struct {
	ResponseCommon
	Data ManyEventSubSubscriptions
}

// Parameter for filtering subscriptions, currently only the status is filterable
type EventSubSubscriptionsParams struct {
	Status string `query:"status"`
	Type   string `query:"type"`
	UserID string `query:"user_id"`
	After  string `query:"after"`
}

// Parameter for removing a subscription.
type RemoveEventSubSubscriptionParams struct {
	ID string `query:"id"`
}

// Response for removing a subscription
type RemoveEventSubSubscriptionParamsResponse struct {
	ResponseCommon
}

// EventSub helper Variables for Types and Status
const (
	EventSubStatusEnabled                      = "enabled"
	EventSubStatusPending                      = "webhook_callback_verification_pending"
	EventSubStatusFailed                       = "webhook_callback_verification_failed"
	EventSubStatusNotificationFailuresExceeded = "notification_failures_exceeded"
	EventSubStatusAuthorizationRevoked         = "authorization_revoked"
	EventSubStatusUserRemoved                  = "user_removed"

	EventSubTypeChannelGoalBegin                          = "channel.goal.begin"
	EventSubTypeChannelGoalProgress                       = "channel.goal.progress"
	EventSubTypeChannelGoalEnd                            = "channel.goal.end"
	EventSubTypeChannelUpdate                             = "channel.update"
	EventSubTypeChannelFollow                             = "channel.follow"
	EventSubTypeChannelSubscription                       = "channel.subscribe"
	EventSubTypeChannelSubscriptionEnd                    = "channel.subscription.end"
	EventSubTypeChannelSubscriptionGift                   = "channel.subscription.gift"
	EventSubTypeChannelSubscriptionMessage                = "channel.subscription.message"
	EventSubTypeChannelCheer                              = "channel.cheer"
	EventSubTypeChannelRaid                               = "channel.raid"
	EventSubTypeChannelBan                                = "channel.ban"
	EventSubTypeChannelUnban                              = "channel.unban"
	EventSubTypeModeratorAdd                              = "channel.moderator.add"
	EventSubTypeModeratorRemove                           = "channel.moderator.remove"
	EventSubTypeChannelPointsCustomRewardAdd              = "channel.channel_points_custom_reward.add"
	EventSubTypeChannelPointsCustomRewardUpdate           = "channel.channel_points_custom_reward.update"
	EventSubTypeChannelPointsCustomRewardRemove           = "channel.channel_points_custom_reward.remove"
	EventSubTypeChannelPointsCustomRewardRedemptionAdd    = "channel.channel_points_custom_reward_redemption.add"
	EventSubTypeChannelPointsCustomRewardRedemptionUpdate = "channel.channel_points_custom_reward_redemption.update"
	EventSubTypeChannelChatClear                          = "channel.chat.clear"
	EventSubTypeChannelChatClearUserMessages              = "channel.chat.clear_user_messages"
	EventSubTypeChannelChatMessage                        = "channel.chat.message"
	EventSubTypeChannelChatMessageDelete                  = "channel.chat.message_delete"
	EventSubTypeChannelChatNotification                   = "channel.chat.notification"
	EventSubTypeChannelPollBegin                          = "channel.poll.begin"
	EventSubTypeChannelPollProgress                       = "channel.poll.progress"
	EventSubTypeChannelPollEnd                            = "channel.poll.end"
	EventSubTypeChannelPredictionBegin                    = "channel.prediction.begin"
	EventSubTypeChannelPredictionProgress                 = "channel.prediction.progress"
	EventSubTypeChannelPredictionLock                     = "channel.prediction.lock"
	EventSubTypeChannelPredictionEnd                      = "channel.prediction.end"
	EventSubExtensionBitsTransactionCreate                = "extension.bits_transaction.create"
	EventSubTypeHypeTrainBegin                            = "channel.hype_train.begin"
	EventSubTypeHypeTrainProgress                         = "channel.hype_train.progress"
	EventSubTypeHypeTrainEnd                              = "channel.hype_train.end"
	EventSubTypeCharityDonation                           = "channel.charity_campaign.donate"
	EventSubTypeCharityProgress                           = "channel.charity_campaign.progress"
	EventSubTypeCharityStop                               = "channel.charity_campaign.stop"
	EventSubTypeCharityStart                              = "channel.charity_campaign.start"
	EventSubTypeStreamOnline                              = "stream.online"
	EventSubTypeStreamOffline                             = "stream.offline"
	EventSubTypeUserAuthorizationRevoke                   = "user.authorization.revoke"
	EventSubTypeUserUpdate                                = "user.update"
	EventSubShoutoutCreate                                = "channel.shoutout.create"
	EventSubShoutoutReceive                               = "channel.shoutout.receive"
)

// Event Notification Responses

// Data for a channel ban notification
type EventSubChannelBanEvent struct {
	UserID               string `json:"user_id"`
	UserLogin            string `json:"user_login"`
	UserName             string `json:"user_name"`
	BroadcasterUserID    string `json:"broadcaster_user_id"`
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	BroadcasterUserName  string `json:"broadcaster_user_name"`
	ModeratorUserID      string `json:"moderator_user_id"`
	ModeratorUserLogin   string `json:"moderator_user_login"`
	ModeratorUserName    string `json:"moderator_user_name"`
	Reason               string `json:"reason"`
	EndsAt               Time   `json:"ends_at"`
	IsPermanent          bool   `json:"is_permanent"`
}

// Data for a channel subscribe notification
type EventSubChannelSubscribeEvent struct {
	UserID               string `json:"user_id"`
	UserLogin            string `json:"user_login"`
	UserName             string `json:"user_name"`
	BroadcasterUserID    string `json:"broadcaster_user_id"`
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	BroadcasterUserName  string `json:"broadcaster_user_name"`
	Tier                 string `json:"tier"`
	IsGift               bool   `json:"is_gift"`
}

// EventSubChannelSubscriptionGiftEvent
type EventSubChannelSubscriptionGiftEvent struct {
	UserID               string `json:"user_id"`
	UserLogin            string `json:"user_login"`
	UserName             string `json:"user_name"`
	BroadcasterUserID    string `json:"broadcaster_user_id"`
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	BroadcasterUserName  string `json:"broadcaster_user_name"`
	Total                int    `json:"total"`
	Tier                 string `json:"tier"`
	CumulativeTotal      int    `json:"cumulative_total"`
	IsAnonymous          bool   `json:"is_anonymous"`
}

// EventSubChannelSubscriptionMessageEvent
type EventSubChannelSubscriptionMessageEvent struct {
	UserID               string          `json:"user_id"`
	UserLogin            string          `json:"user_login"`
	UserName             string          `json:"user_name"`
	BroadcasterUserID    string          `json:"broadcaster_user_id"`
	BroadcasterUserLogin string          `json:"broadcaster_user_login"`
	BroadcasterUserName  string          `json:"broadcaster_user_name"`
	Tier                 string          `json:"tier"`
	Message              EventSubMessage `json:"message"`
	CumulativeMonths     int             `json:"cumulative_months"`
	StreakMonths         int             `json:"streak_months"`
	DurationMonths       int             `json:"duration_months"`
}

// Data for a channel cheer notification
type EventSubChannelCheerEvent struct {
	IsAnonymous          bool   `json:"is_anonymous"`
	UserID               string `json:"user_id"`
	UserLogin            string `json:"user_login"`
	UserName             string `json:"user_name"`
	BroadcasterUserID    string `json:"broadcaster_user_id"`
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	BroadcasterUserName  string `json:"broadcaster_user_name"`
	Message              string `json:"message"`
	Bits                 int    `json:"bits"`
}

// Data for a channel update notification
type EventSubChannelUpdateEvent struct {
	BroadcasterUserID    string `json:"broadcaster_user_id"`
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	BroadcasterUserName  string `json:"broadcaster_user_name"`
	Title                string `json:"title"`
	Language             string `json:"language"`
	CategoryID           string `json:"category_id"`
	CategoryName         string `json:"category_name"`
	IsMature             bool   `json:"is_mature"`
}

// Data for a channel unban notification
type EventSubChannelUnbanEvent struct {
	UserID               string `json:"user_id"`
	UserLogin            string `json:"user_login"`
	UserName             string `json:"user_name"`
	BroadcasterUserID    string `json:"broadcaster_user_id"`
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	BroadcasterUserName  string `json:"broadcaster_user_name"`
	ModeratorUserID      string `json:"moderator_user_id"`
	ModeratorUserLogin   string `json:"moderator_user_login"`
	ModeratorUserName    string `json:"moderator_user_name"`
}

// Data for a channel follow notification
type EventSubChannelFollowEvent struct {
	UserID               string `json:"user_id"`
	UserLogin            string `json:"user_login"`
	UserName             string `json:"user_name"`
	BroadcasterUserID    string `json:"broadcaster_user_id"`
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	BroadcasterUserName  string `json:"broadcaster_user_name"`
	FollowedAt           Time   `json:"followed_at"`
}

// Data for a channel moderator add notification, it's the same as the channel follow notification
type EventSubModeratorAddEvent = EventSubChannelFollowEvent

// Data for a channel moderator remove notification, it's the same as the channel follow notification
type EventSubModeratorRemoveEvent = EventSubChannelFollowEvent

// Data for a channel raid notification
type EventSubChannelRaidEvent struct {
	FromBroadcasterUserID    string `json:"from_broadcaster_user_id"`
	FromBroadcasterUserLogin string `json:"from_broadcaster_user_login"`
	FromBroadcasterUserName  string `json:"from_broadcaster_user_name"`
	ToBroadcasterUserID      string `json:"to_broadcaster_user_id"`
	ToBroadcasterUserLogin   string `json:"to_broadcaster_user_login"`
	ToBroadcasterUserName    string `json:"to_broadcaster_user_name"`
	Viewers                  int    `json:"viewers"`
}

// Data for a chat clear event
type EventSubChannelChatClearEvent struct {
	BroadcasterUserID    string `json:"broadcaster_user_id"`
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	BroadcasterUserName  string `json:"broadcaster_user_name"`
}

// Data for a chat clear user messages event
type EventSubChannelChatClearUserMessagesEvent struct {
	BroadcasterUserID    string `json:"broadcaster_user_id"`
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	BroadcasterUserName  string `json:"broadcaster_user_name"`
	TargetUserID         string `json:"target_user_id"`
	TargetUserLogin      string `json:"target_user_login"`
	TargetUserName       string `json:"target_user_name"`
}

// Data for a chat message
type EventSubChannelChatMessageEvent struct {
	BroadcasterUserID           string                   `json:"broadcaster_user_id"`
	BroadcasterUserLogin        string                   `json:"broadcaster_user_login"`
	BroadcasterUserName         string                   `json:"broadcaster_user_name"`
	ChatterUserID               string                   `json:"chatter_user_id"`
	ChatterUserLogin            string                   `json:"chatter_user_login"`
	ChatterUserName             string                   `json:"chatter_user_name"`
	MessageID                   string                   `json:"message_id"`
	Message                     EventSubChatMessage      `json:"message"`
	MessageType                 EventSubChatMessageType  `json:"message_type"`
	Badges                      []EventSubChatBadge      `json:"badges"`
	Cheer                       EventSubChatMessageCheer `json:"cheer"`
	Color                       string                   `json:"color"`
	Reply                       EventSubChatMessageReply `json:"reply"`
	ChannelPointsCustomRewardID string                   `json:"channel_points_custom_reward_id"`
}

type EventSubChatMessage struct {
	Text      string                        `json:"text"`
	Fragments []EventSubChatMessageFragment `json:"fragments"`
}

type EventSubChatMessageReply struct {
	ParentMessageID   string `json:"parent_message_id"`
	ParentMessageBody string `json:"parent_message_body"`
	ParentUserID      string `json:"parent_user_id"`
	ParentUserName    string `json:"parent_user_name"`
	ParentUserLogin   string `json:"parent_user_login"`
	ThreadMessageID   string `json:"thread_message_id"`
	ThreadUserID      string `json:"thread_user_id"`
	ThreadUserName    string `json:"thread_user_name"`
	ThreadUserLogin   string `json:"thread_user_login"`
}

type EventSubChatMessageCheer struct {
	Bits int64 `json:"bits"`
}

type EventSubChatBadge struct {
	SetID string `json:"set_id"`
	ID    string `json:"id"`
	Info  string `json:"info"`
}

type EventSubChatMessageType string

const (
	EventSubChatMessageTypeText                     EventSubChatMessageType = "text"
	EventSubChatMessageTypeChannelPointsHighlighted EventSubChatMessageType = "channel_points_highlighted"
	EventSubChatMessageTypeChannelPointsSubOnly     EventSubChatMessageType = "channel_points_sub_only"
	EventSubChatMessageTypeUserIntro                EventSubChatMessageType = "user_intro"
)

type EventSubChatMessageFragmentType string

const (
	EventSubChatMessageFragmentTypeText      EventSubChatMessageFragmentType = "text"
	EventSubChatMessageFragmentTypeCheermote EventSubChatMessageFragmentType = "cheermote"
	EventSubChatMessageFragmentTypeEmote     EventSubChatMessageFragmentType = "emote"
	EventSubChatMessageFragmentTypeMention   EventSubChatMessageFragmentType = "mention"
)

type EventSubChatMessageFragment struct {
	Type      EventSubChatMessageFragmentType `json:"type"`
	Text      string                          `json:"text"`
	Cheermote EventSubChatMessageCheermote    `json:"cheermote"`
	Emote     EventSubChatMessageEmote        `json:"emote"`
	Mention   EventSubChatMessageMention      `json:"mention"`
}

type EventSubChatMessageCheermote struct {
	Prefix string `json:"prefix"`
	Bits   int64  `json:"bits"`
	Tier   int    `json:"tier"`
}

type EventSubChatMessageEmote struct {
	ID         string   `json:"id"`
	EmoteSetID string   `json:"emote_set_id"`
	OwnerID    string   `json:"owner_id"`
	Format     []string `json:"format"`
}

type EventSubChatMessageMention struct {
	UserID    string `json:"user_id"`
	UserName  string `json:"user_name"`
	UserLogin string `json:"user_login"`
}

// Data for a chat message delete event
type EventSubChannelChatMessageDeleteEvent struct {
	BroadcasterUserID    string `json:"broadcaster_user_id"`
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	BroadcasterUserName  string `json:"broadcaster_user_name"`
	TargetUserID         string `json:"target_user_id"`
	TargetUserLogin      string `json:"target_user_login"`
	TargetUserName       string `json:"target_user_name"`
	MessageID            string `json:"message_id"`
}

// Data for a chat notification event
type EventSubChannelChatNotificationEvent struct {
	BroadcasterUserID    string                                          `json:"broadcaster_user_id"`
	BroadcasterUserLogin string                                          `json:"broadcaster_user_login"`
	BroadcasterUserName  string                                          `json:"broadcaster_user_name"`
	ChatterUserID        string                                          `json:"chatter_user_id"`
	ChatterUserLogin     string                                          `json:"chatter_user_login"`
	ChatterUserName      string                                          `json:"chatter_user_name"`
	ChatterIsAnonymous   bool                                            `json:"chatter_is_anonymous"`
	Color                string                                          `json:"color"`
	Badges               []EventSubChatBadge                             `json:"badges"`
	SystemMessage        string                                          `json:"system_message"`
	MessageID            string                                          `json:"message_id"`
	Message              EventSubChatNotificationMessage                 `json:"message"`
	NoticeType           EventSubChannelChatNotificationType             `json:"notice_type"`
	Sub                  EventSubChannelChatNotificationSub              `json:"sub"`
	Resub                EventSubChannelChatNotificationResub            `json:"resub"`
	SubGift              EventSubChannelChatNotificationSubGift          `json:"sub_gift"`
	CommunitySubGift     EventSubChannelChatNotificationCommunitySubGift `json:"community_sub_gift"`
	GiftPaidUpgrade      EventSubChannelChatNotificationGiftPaidUpgrade  `json:"gift_paid_upgrade"`
	PrimePaidUpgrade     EventSubChannelChatNotificationPrimePaidUpgrade `json:"prime_paid_upgrade"`
	Raid                 EventSubChannelChatNotificationRaid             `json:"raid"`
	Unraid               EventSubChannelChatNotificationUnraid           `json:"unraid"`
	PayItForward         EventSubChannelChatNotificationPayItForward     `json:"pay_it_forward"`
	Announcement         EventSubChannelChatNotificationAnnouncement     `json:"announcement"`
	CharityDonation      EventSubChannelChatNotificationCharityDonation  `json:"charity_donation"`
	BitsBadgeTier        EventSubChannelChatNotificationBitsBadgeTier    `json:"bits_badge_tier"`
}

type EventSubChannelChatNotificationType string

const (
	EventSubChannelNotificationSub              EventSubChannelChatNotificationType = "sub"
	EventSubChannelNotificationResub            EventSubChannelChatNotificationType = "resub"
	EventSubChannelNotificationSubGift          EventSubChannelChatNotificationType = "sub_gift"
	EventSubChannelNotificationCommunitySubGift EventSubChannelChatNotificationType = "community_sub_gift"
	EventSubChannelNotificationGiftPaidUpgrade  EventSubChannelChatNotificationType = "gift_paid_upgrade"
	EventSubChannelNotificationPrimePaidUpgrade EventSubChannelChatNotificationType = "prime_paid_upgrade"
	EventSubChannelNotificationRaid             EventSubChannelChatNotificationType = "raid"
	EventSubChannelNotificationUnraid           EventSubChannelChatNotificationType = "unraid"
	EventSubChannelNotificationPayItForward     EventSubChannelChatNotificationType = "pay_it_forward"
	EventSubChannelNotificationAnnouncement     EventSubChannelChatNotificationType = "announcement"
	EventSubChannelNotificationBitsBadgeTier    EventSubChannelChatNotificationType = "bits_badge_tier"
	EventSubChannelNotificationCharityDonation  EventSubChannelChatNotificationType = "charity_donation"
)

type EventSubChannelChatNotificationSub struct {
	SubTier        string `json:"sub_tier"`
	IsPrime        bool   `json:"is_prime"`
	DurationMonths int    `json:"duration_months"`
}

type EventSubChannelChatNotificationResub struct {
	CumulativeMonths  int    `json:"cumulative_months"`
	DurationMonths    int    `json:"duration_months"`
	StreakMonths      int    `json:"streak_months"`
	SubTier           string `json:"sub_tier"`
	IsPrime           bool   `json:"is_prime"`
	IsGift            bool   `json:"is_gift"`
	GifterIsAnonymous bool   `json:"gifter_is_anonymous"`
	GifterUserID      string `json:"gifter_user_id"`
	GifterUserLogin   string `json:"gifter_user_login"`
	GifterUserName    string `json:"gifter_user_name"`
}

type EventSubChannelChatNotificationSubGift struct {
	DurationMonths     int    `json:"duration_months"`
	CumulativeTotal    int    `json:"cumulative_total"`
	RecipientUserID    string `json:"recipient_user_id"`
	RecipientUserLogin string `json:"recipient_user_login"`
	RecipientUserName  string `json:"recipient_user_name"`
	SubTier            string `json:"sub_tier"`
	CommunityGiftID    string `json:"community_gift_id"`
}

type EventSubChannelChatNotificationCommunitySubGift struct {
	ID              string `json:"id"`
	Total           int    `json:"total"`
	SubTier         string `json:"sub_tier"`
	CumulativeTotal int    `json:"cumulative_total"`
}

type EventSubChannelChatNotificationGiftPaidUpgrade struct {
	GifterIsAnonymous bool   `json:"gifter_is_anonymous"`
	GifterUserID      string `json:"gifter_user_id"`
	GifterUserLogin   string `json:"gifter_user_login"`
	GifterUserName    string `json:"gifter_user_name"`
}

type EventSubChannelChatNotificationPrimePaidUpgrade struct {
	SubTier string `json:"sub_tier"`
}

type EventSubChannelChatNotificationRaid struct {
	UserID          string `json:"user_id"`
	UserLogin       string `json:"user_login"`
	UserName        string `json:"user_name"`
	ViewerCount     int64  `json:"viewer_count"`
	ProfileImageURL string `json:"profile_image_url"`
}

type EventSubChannelChatNotificationUnraid struct{}

type EventSubChannelChatNotificationPayItForward struct {
	GifterIsAnonymous bool   `json:"gifter_is_anonymous"`
	GifterUserID      string `json:"gifter_user_id"`
	GifterUserLogin   string `json:"gifter_user_login"`
	GifterUserName    string `json:"gifter_user_name"`
}

type EventSubChannelChatNotificationAnnouncement struct {
	Color string `json:"color"`
}

type EventSubChannelChatNotificationCharityDonation struct {
	CharityName string                                               `json:"charity_name"`
	Amount      EventSubChannelChatNotificationCharityDonationAmount `json:"amount"`
}

type EventSubChannelChatNotificationCharityDonationAmount struct {
	Value        int64  `json:"value"`
	DecimalPlace int64  `json:"decimal_place"`
	Currency     string `json:"currency"`
}

type EventSubChannelChatNotificationBitsBadgeTier struct {
	Tier int64 `json:"tier"`
}

type EventSubChatNotificationMessage struct {
	Text      string                        `json:"text"`
	Fragments []EventSubChatMessageFragment `json:"fragments"`
}

// Data for a channel poll begin event
type EventSubChannelPollBeginEvent struct {
	ID                   string                      `json:"id"`
	BroadcasterUserID    string                      `json:"broadcaster_user_id"`
	BroadcasterUserLogin string                      `json:"broadcaster_user_login"`
	BroadcasterUserName  string                      `json:"broadcaster_user_name"`
	Title                string                      `json:"title"`
	Choices              []PollChoice                `json:"choices"`
	BitsVoting           EventSubBitVoting           `json:"bits_voting"`
	ChannelPointsVoting  EventSubChannelPointsVoting `json:"channel_points_voting"`
	StartedAt            Time                        `json:"started_at"`
	EndsAt               Time                        `json:"ends_at"`
}

// Data for a channel poll progress event, it's the same as the channel poll begin event
type EventSubChannelPollProgressEvent = EventSubChannelPollBeginEvent

// Data for a channel poll end event
type EventSubChannelPollEndEvent struct {
	ID                   string                      `json:"id"`
	BroadcasterUserID    string                      `json:"broadcaster_user_id"`
	BroadcasterUserLogin string                      `json:"broadcaster_user_login"`
	BroadcasterUserName  string                      `json:"broadcaster_user_name"`
	Title                string                      `json:"title"`
	Choices              []PollChoice                `json:"choices"`
	BitsVoting           EventSubBitVoting           `json:"bits_voting"`
	ChannelPointsVoting  EventSubChannelPointsVoting `json:"channel_points_voting"`
	Status               string                      `json:"status"`
	StartedAt            Time                        `json:"started_at"`
	EndedAt              Time                        `json:"ended_at"`
}

type EventSubBitVoting struct {
	IsEnabled     bool `json:"is_enabled"`
	AmountPerVote int  `json:"amount_per_vote"`
}

type EventSubChannelPointsVoting = EventSubBitVoting

// Data for a channel points custom reward notification
type EventSubChannelPointsCustomRewardEvent struct {
	ID                                string                 `json:"id"`
	BroadcasterUserID                 string                 `json:"broadcaster_user_id"`
	BroadcasterUserLogin              string                 `json:"broadcaster_user_login"`
	BroadcasterUserName               string                 `json:"broadcaster_user_name"`
	IsEnabled                         bool                   `json:"is_enabled"`
	IsPaused                          bool                   `json:"is_paused"`
	IsInStock                         bool                   `json:"is_in_stock"`
	Title                             string                 `json:"title"`
	Cost                              int                    `json:"cost"`
	Prompt                            string                 `json:"prompt"`
	IsUserInputRequired               bool                   `json:"is_user_input_required"`
	ShouldRedemptionsSkipRequestQueue bool                   `json:"should_redemptions_skip_request_queue"`
	MaxPerStream                      EventSubMaxPerStream   `json:"max_per_stream"`
	MaxPerUserPerStream               EventSubMaxPerStream   `json:"max_per_user_per_stream"`
	BackgroundColor                   string                 `json:"background_color"`
	Image                             EventSubImage          `json:"image"`
	DefaultImage                      EventSubImage          `json:"default_image"`
	GlobalCooldown                    EventSubGlobalCooldown `json:"global_cooldown"`
	CooldownExpiresAt                 Time                   `json:"cooldown_expires_at"`
	RedemptionsRedeemedCurrentStream  int                    `json:"redemptions_redeemed_current_stream"`
}

// Data for a channel points custom reward redemption notification
type EventSubChannelPointsCustomRewardRedemptionEvent struct {
	ID                   string         `json:"id"`
	BroadcasterUserID    string         `json:"broadcaster_user_id"`
	BroadcasterUserLogin string         `json:"broadcaster_user_login"`
	BroadcasterUserName  string         `json:"broadcaster_user_name"`
	UserID               string         `json:"user_id"`
	UserLogin            string         `json:"user_login"`
	UserName             string         `json:"user_name"`
	UserInput            string         `json:"user_input"`
	Status               string         `json:"status"`
	Reward               EventSubReward `json:"reward"`
	RedeemedAt           Time           `json:"redeemed_at"`
}

// Data for a channel prediction begin event
type EventSubChannelPredictionBeginEvent struct {
	ID                   string            `json:"id"`
	BroadcasterUserID    string            `json:"broadcaster_user_id"`
	BroadcasterUserLogin string            `json:"broadcaster_user_login"`
	BroadcasterUserName  string            `json:"broadcaster_user_name"`
	Title                string            `json:"title"`
	Outcomes             []EventSubOutcome `json:"outcomes"`
	StartedAt            Time              `json:"started_at"`
	LocksAt              Time              `json:"locks_at"`
}

// Data for a channel prediction progress event
type EventSubChannelPredictionProgressEvent = EventSubChannelPredictionBeginEvent

// Data for a channel prediction lock event
type EventSubChannelPredictionLockEvent struct {
	ID                   string            `json:"id"`
	BroadcasterUserID    string            `json:"broadcaster_user_id"`
	BroadcasterUserLogin string            `json:"broadcaster_user_login"`
	BroadcasterUserName  string            `json:"broadcaster_user_name"`
	Title                string            `json:"title"`
	WinningOutcomeID     string            `json:"winning_outcome_id"`
	Outcomes             []EventSubOutcome `json:"outcomes"`
	Status               string            `json:"status"`
	StartedAt            Time              `json:"started_at"`
	LockedAt             Time              `json:"locked_at"`
}

// Data for a channel prediction end event
type EventSubChannelPredictionEndEvent struct {
	ID                   string            `json:"id"`
	BroadcasterUserID    string            `json:"broadcaster_user_id"`
	BroadcasterUserLogin string            `json:"broadcaster_user_login"`
	BroadcasterUserName  string            `json:"broadcaster_user_name"`
	Title                string            `json:"title"`
	WinningOutcomeID     string            `json:"winning_outcome_id"`
	Outcomes             []EventSubOutcome `json:"outcomes"`
	Status               string            `json:"status"`
	StartedAt            Time              `json:"started_at"`
	EndedAt              Time              `json:"ended_at"`
}

// Data for an extension bits transaction creation
type EventSubExtensionBitsTransactionCreateEvent struct {
	ExtensionClientID    string          `json:"extension_client_id"`
	ID                   string          `json:"id"`
	BroadcasterUserID    string          `json:"broadcaster_user_id"`
	BroadcasterUserLogin string          `json:"broadcaster_user_login"`
	BroadcasterUserName  string          `json:"broadcaster_user_name"`
	UserID               string          `json:"user_id"`
	UserLogin            string          `json:"user_login"`
	UserName             string          `json:"user_name"`
	Product              EventSubProduct `json:"product"`
}

// Data for a hype train begin notification
type EventSubHypeTrainBeginEvent struct {
	BroadcasterUserID    string                 `json:"broadcaster_user_id"`
	BroadcasterUserLogin string                 `json:"broadcaster_user_login"`
	BroadcasterUserName  string                 `json:"broadcaster_user_name"`
	Total                int                    `json:"total"`
	Progress             int                    `json:"progress"`
	Goal                 int                    `json:"goal"`
	TopContributions     []EventSubContribution `json:"top_contributions"`
	LastContribution     EventSubContribution   `json:"last_contribution"`
	StartedAt            Time                   `json:"started_at"`
	ExpiresAt            Time                   `json:"expires_at"`
	IsGoldenKappaTrain   bool                   `json:"is_golden_kappa_train"`
}

// Data for a hype train progress notification
type EventSubHypeTrainProgressEvent struct {
	BroadcasterUserID    string                 `json:"broadcaster_user_id"`
	BroadcasterUserLogin string                 `json:"broadcaster_user_login"`
	BroadcasterUserName  string                 `json:"broadcaster_user_name"`
	Level                int                    `json:"level"`
	Total                int                    `json:"total"`
	Progress             int                    `json:"progress"`
	Goal                 int                    `json:"goal"`
	TopContributions     []EventSubContribution `json:"top_contributions"`
	LastContribution     EventSubContribution   `json:"last_contribution"`
	StartedAt            Time                   `json:"started_at"`
	ExpiresAt            Time                   `json:"expires_at"`
	IsGoldenKappaTrain   bool                   `json:"is_golden_kappa_train"`
}

// Data for a hype train end notification
type EventSubHypeTrainEndEvent struct {
	ID                   string                 `json:"id"`
	BroadcasterUserID    string                 `json:"broadcaster_user_id"`
	BroadcasterUserLogin string                 `json:"broadcaster_user_login"`
	BroadcasterUserName  string                 `json:"broadcaster_user_name"`
	Level                int                    `json:"level"`
	Total                int                    `json:"total"`
	TopContributions     []EventSubContribution `json:"top_contributions"`
	StartedAt            Time                   `json:"started_at"`
	EndedAt              Time                   `json:"ended_at"`
	CooldownEndsAt       Time                   `json:"cooldown_ends_at"`
	IsGoldenKappaTrain   bool                   `json:"is_golden_kappa_train"`
}

// Data for a stream online notification
type EventSubStreamOnlineEvent struct {
	ID                   string `json:"id"`
	BroadcasterUserID    string `json:"broadcaster_user_id"`
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	BroadcasterUserName  string `json:"broadcaster_user_name"`
	Type                 string `json:"type"`
	StartedAt            Time   `json:"started_at"`
}

// Data for a stream offline notification
type EventSubStreamOfflineEvent struct {
	BroadcasterUserID    string `json:"broadcaster_user_id"`
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	BroadcasterUserName  string `json:"broadcaster_user_name"`
}

// Data for an user authentication revoke notification, this means the user has revoked the access token and if you need to comply with gdpr you need to delete your user data belonging to the user.
type EventSubUserAuthenticationRevokeEvent struct {
	ClientID  string `json:"client_id"`
	UserID    string `json:"user_id"`
	UserLogin string `json:"user_login"`
	UserName  string `json:"user_name"`
}

// Data for an user update notification
type EventSubUserUpdateEvent struct {
	UserID      string `json:"user_id"`
	UserLogin   string `json:"user_login"`
	UserName    string `json:"user_name"`
	Email       string `json:"email"`
	Description string `json:"description"`
}

// This belongs to a custom reward and defines it's cooldown
type EventSubGlobalCooldown struct {
	IsEnabled bool `json:"is_enabled"`
	Seconds   int  `json:"seconds"`
}

// This also belongs to a custom reward and defines the image urls
type EventSubImage struct {
	Url1x string `json:"url_1x"`
	Url2x string `json:"url_2x"`
	Url4x string `json:"url_4x"`
}

// This belongs to a hype train and defines a user contribution
type EventSubContribution struct {
	UserID    string `json:"user_id"`
	UserLogin string `json:"user_login"`
	UserName  string `json:"user_name"`
	Type      string `json:"type"`
	Total     int64  `json:"total"`
}

// This belong to an outcome and defines user reward
type EventSubTopPredictor struct {
	UserID            string `json:"user_id"`
	UserLogin         string `json:"user_login"`
	UserName          string `json:"user_name"`
	ChannelPointWon   int    `json:"channel_points_won"`
	ChannelPointsUsed int    `json:"channel_points_used"`
}

// This belongs to a custom reward and defines if it is limited per stream
type EventSubMaxPerStream struct {
	IsEnabled bool `json:"is_enabled"`
	Value     int  `json:"value"`
}

// This belong to a channel prediction and defines the outcomes
type EventSubOutcome struct {
	ID            string                 `json:"id"`
	Title         string                 `json:"title"`
	Color         string                 `json:"color"`
	Users         int                    `json:"users"`
	ChannelPoints int                    `json:"channel_points"`
	TopPredictors []EventSubTopPredictor `json:"top_predictors"`
}

type EventSubProduct struct {
	Name          string `json:"name"`
	Bits          int    `json:"bits"`
	Sku           string `json:"sku"`
	InDevelopment bool   `json:"in_development"`
}

// This belongs to a reward redemption and defines the reward redeemed
type EventSubReward struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Cost   int    `json:"cost"`
	Prompt string `json:"prompt"`
}

// EventSubMessage
type EventSubMessage struct {
	Text   string          `json:"text"`
	Emotes []EventSubEmote `json:"emotes"`
}

// EventSubEmote
type EventSubEmote struct {
	Begin int    `json:"begin"`
	End   int    `json:"end"`
	ID    string `json:"id"`
}

type EventSubChannelGoalStartEvent struct {
	ID                   string `json:"id"`
	BroadcasterUserID    string `json:"broadcaster_user_id"`
	BroadcasterUserName  string `json:"broadcaster_user_name"`
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	Type                 string `json:"type"`
	Description          string `json:"description"`
	CurrentAmount        int    `json:"current_amount"`
	TargetAmount         int    `json:"target_amount"`
	StartedAt            Time   `json:"started_at"`
}

type EventSubChannelGoalProgressEvent struct {
	ID                   string `json:"id"`
	BroadcasterUserID    string `json:"broadcaster_user_id"`
	BroadcasterUserName  string `json:"broadcaster_user_name"`
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	Type                 string `json:"type"`
	Description          string `json:"description"`
	CurrentAmount        int    `json:"current_amount"`
	TargetAmount         int    `json:"target_amount"`
	StartedAt            Time   `json:"started_at"`
}

type EventSubChannelGoalEndEvent struct {
	ID                   string `json:"id"`
	BroadcasterUserID    string `json:"broadcaster_user_id"`
	BroadcasterUserName  string `json:"broadcaster_user_name"`
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	Type                 string `json:"type"`
	Description          string `json:"description"`
	IsAchieved           bool   `json:"is_achieved"`
	CurrentAmount        int    `json:"current_amount"`
	TargetAmount         int    `json:"target_amount"`
	StartedAt            Time   `json:"started_at"`
	EndedAt              Time   `json:"ended_at"`
}

type EventSubCharityAmount struct {
	Value         int64  `json:"value"`
	DecimalPlaces int64  `json:"decimal_places"`
	Currency      string `json:"currency"`
}

type EventSubCharityDonationEvent struct {
	DonationID           string                `json:"id"`
	CharityCampaignID    string                `json:"campaign_id"`
	CharityDescription   string                `json:"campaign_description"`
	CharityWebsite       string                `json:"campaign_website"`
	CharityName          string                `json:"charity_name"`
	CharityLogoURL       string                `json:"charity_logo"`
	BroadcasterUserID    string                `json:"broadcaster_user_id"`
	BroadcasterUserName  string                `json:"broadcaster_user_name"`
	BroadcasterUserLogin string                `json:"broadcaster_user_login"`
	UserID               string                `json:"user_id"`
	UserName             string                `json:"user_name"`
	UserLogin            string                `json:"user_login"`
	Amount               EventSubCharityAmount `json:"amount"`
}

type EventSubCharityProgressEvent struct {
	CharityCampaignID    string                `json:"campaign_id"`
	CharityDescription   string                `json:"campaign_description"`
	CharityWebsite       string                `json:"campaign_website"`
	CharityName          string                `json:"charity_name"`
	CharityLogoURL       string                `json:"charity_logo"`
	BroadcasterUserID    string                `json:"broadcaster_id"`
	BroadcasterUserName  string                `json:"broadcaster_name"`
	BroadcasterUserLogin string                `json:"broadcaster_user_login"`
	UserID               string                `json:"user_id"`
	UserName             string                `json:"user_name"`
	UserLogin            string                `json:"user_login"`
	Amount               EventSubCharityAmount `json:"amount"`
}

type EventSubCharityStopEvent struct {
	CharityCampaignID    string                `json:"campaign_id"`
	CharityDescription   string                `json:"campaign_description"`
	CharityWebsite       string                `json:"campaign_website"`
	CharityName          string                `json:"charity_name"`
	CharityLogoURL       string                `json:"charity_logo"`
	BroadcasterUserID    string                `json:"broadcaster_id"`
	BroadcasterUserName  string                `json:"broadcaster_name"`
	BroadcasterUserLogin string                `json:"broadcaster_login"`
	UserID               string                `json:"user_id"`
	UserName             string                `json:"user_name"`
	UserLogin            string                `json:"user_login"`
	CurrentAmount        EventSubCharityAmount `json:"current_amount"`
	TargetAmount         EventSubCharityAmount `json:"target_amount"`
	StoppedAt            Time                  `json:"stopped_at"`
}

type EventSubCharityStartEvent struct {
	CharityCampaignID    string                `json:"campaign_id"`
	CharityDescription   string                `json:"campaign_description"`
	CharityWebsite       string                `json:"campaign_website"`
	CharityName          string                `json:"charity_name"`
	CharityLogoURL       string                `json:"charity_logo"`
	BroadcasterUserID    string                `json:"broadcaster_id"`
	BroadcasterUserName  string                `json:"broadcaster_name"`
	BroadcasterUserLogin string                `json:"broadcaster_login"`
	UserID               string                `json:"user_id"`
	UserName             string                `json:"user_name"`
	UserLogin            string                `json:"user_login"`
	CurrentAmount        EventSubCharityAmount `json:"current_amount"`
	TargetAmount         EventSubCharityAmount `json:"target_amount"`
	StartedAt            Time                  `json:"started_at"`
}

type EventSubShoutoutCreateEvent struct {
	BroadcasterUserID      string `json:"broadcaster_user_id"`
	BroadcasterUserName    string `json:"broadcaster_user_name"`
	BroadcasterUserLogin   string `json:"broadcaster_user_login"`
	ModeratorUserID        string `json:"moderator_user_id"`
	ModeratorUserName      string `json:"moderator_user_name"`
	ModeratorUserLogin     string `json:"moderator_user_login"`
	ToBroadcasterUserID    string `json:"to_broadcaster_user_id"`
	ToBroadcasterUserName  string `json:"to_broadcaster_user_name"`
	ToBroadcasterUserLogin string `json:"to_broadcaster_user_login"`
	StartedAt              Time   `json:"started_at"`
	ViewerCount            int64  `json:"viewer_count"`
	CooldownEndsAt         Time   `json:"cooldown_ends_at"`
	TargetCooldownEndsAt   Time   `json:"target_cooldown_ends_at"`
}

type EventSubShoutoutReceiveEvent struct {
	BroadcasterUserID        string `json:"broadcaster_user_id"`
	BroadcasterUserName      string `json:"broadcaster_user_name"`
	BroadcasterUserLogin     string `json:"broadcaster_user_login"`
	FromBroadcasterUserID    string `json:"from_broadcaster_user_id"`
	FromBroadcasterUserName  string `json:"from_broadcaster_user_name"`
	FromBroadcasterUserLogin string `json:"from_broadcaster_user_login"`
	ViewerCount              int64  `json:"viewer_count"`
	StartedAt                Time   `json:"started_at"`
}

// Get all EventSub Subscriptions
func (c *Client) GetEventSubSubscriptions(params *EventSubSubscriptionsParams) (*EventSubSubscriptionsResponse, error) {
	resp, err := c.get("/eventsub/subscriptions", &ManyEventSubSubscriptions{}, params)
	if err != nil {
		return nil, err
	}

	eventSubs := &EventSubSubscriptionsResponse{}
	resp.HydrateResponseCommon(&eventSubs.ResponseCommon)
	eventSubs.Data.Total = resp.Data.(*ManyEventSubSubscriptions).Total
	eventSubs.Data.TotalCost = resp.Data.(*ManyEventSubSubscriptions).TotalCost
	eventSubs.Data.MaxTotalCost = resp.Data.(*ManyEventSubSubscriptions).MaxTotalCost
	eventSubs.Data.EventSubSubscriptions = resp.Data.(*ManyEventSubSubscriptions).EventSubSubscriptions
	eventSubs.Data.Pagination = resp.Data.(*ManyEventSubSubscriptions).Pagination

	return eventSubs, nil
}

// Remove an EventSub Subscription
func (c *Client) RemoveEventSubSubscription(id string) (*RemoveEventSubSubscriptionParamsResponse, error) {

	resp, err := c.delete("/eventsub/subscriptions", nil, &RemoveEventSubSubscriptionParams{ID: id})
	if err != nil {
		return nil, err
	}

	eventsub := &RemoveEventSubSubscriptionParamsResponse{}
	resp.HydrateResponseCommon(&eventsub.ResponseCommon)
	return eventsub, nil
}

// Creates an EventSub subscription
func (c *Client) CreateEventSubSubscription(payload *EventSubSubscription) (*EventSubSubscriptionsResponse, error) {
	switch payload.Transport.Method {
	case "webhook":
		if err := verifyWebhookSub(payload); err != nil {
			return nil, err
		}
	case "websocket":
		if err := verifyWebsocketSub(payload); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("error: unsupported transport method: %s", payload.Transport.Method)
	}

	resp, err := c.postAsJSON("/eventsub/subscriptions", &ManyEventSubSubscriptions{}, payload)
	if err != nil {
		return nil, err
	}

	eventsub := &EventSubSubscriptionsResponse{}
	resp.HydrateResponseCommon(&eventsub.ResponseCommon)
	eventsub.Data = *resp.Data.(*ManyEventSubSubscriptions)
	return eventsub, nil
}

// Verifys that a notification came from twitch using the a signature and the secret used when creating the subscription
func VerifyEventSubNotification(secret string, header http.Header, message string) bool {
	hmacMessage := []byte(fmt.Sprintf("%s%s%s", header.Get("Twitch-Eventsub-Message-Id"), header.Get("Twitch-Eventsub-Message-Timestamp"), message))
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(hmacMessage)
	hmacsha256 := fmt.Sprintf("sha256=%s", hex.EncodeToString(mac.Sum(nil)))
	return hmacsha256 == header.Get("Twitch-Eventsub-Message-Signature")
}

func verifyWebhookSub(payload *EventSubSubscription) error {
	if !strings.HasPrefix(payload.Transport.Callback, "https://") {
		return fmt.Errorf("error: callback must use https")
	}

	if payload.Transport.Secret != "" && (len(payload.Transport.Secret) < 10 || len(payload.Transport.Secret) > 100) {
		return fmt.Errorf("error: secret must be between 10 and 100 characters")
	}

	callbackUrl, err := url.Parse(payload.Transport.Callback)
	if err != nil {
		return err
	}
	if callbackUrl.Port() != "" && callbackUrl.Port() != "443" {
		return fmt.Errorf("error: callback must use port 443")
	}

	return nil
}

func verifyWebsocketSub(payload *EventSubSubscription) error {
	if len(payload.Transport.SessionID) == 0 {
		return fmt.Errorf("error: session ID must be set up")
	}

	return nil
}
