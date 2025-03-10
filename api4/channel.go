// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package api4

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/mattermost/mattermost-server/v6/app"
	"github.com/mattermost/mattermost-server/v6/audit"
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/mattermost/mattermost-server/v6/shared/mlog"
)

func (api *API) InitChannel() {
	api.BaseRoutes.Channels.Handle("", api.ApiSessionRequired(getAllChannels)).Methods("GET")
	api.BaseRoutes.Channels.Handle("", api.ApiSessionRequired(createChannel)).Methods("POST")
	api.BaseRoutes.Channels.Handle("/direct", api.ApiSessionRequired(createDirectChannel)).Methods("POST")
	api.BaseRoutes.Channels.Handle("/search", api.ApiSessionRequiredDisableWhenBusy(searchAllChannels)).Methods("POST")
	api.BaseRoutes.Channels.Handle("/group/search", api.ApiSessionRequiredDisableWhenBusy(searchGroupChannels)).Methods("POST")
	api.BaseRoutes.Channels.Handle("/group", api.ApiSessionRequired(createGroupChannel)).Methods("POST")
	api.BaseRoutes.Channels.Handle("/members/{user_id:[A-Za-z0-9]+}/view", api.ApiSessionRequired(viewChannel)).Methods("POST")
	api.BaseRoutes.Channels.Handle("/{channel_id:[A-Za-z0-9]+}/scheme", api.ApiSessionRequired(updateChannelScheme)).Methods("PUT")

	api.BaseRoutes.ChannelsForTeam.Handle("", api.ApiSessionRequired(getPublicChannelsForTeam)).Methods("GET")
	api.BaseRoutes.ChannelsForTeam.Handle("/deleted", api.ApiSessionRequired(getDeletedChannelsForTeam)).Methods("GET")
	api.BaseRoutes.ChannelsForTeam.Handle("/private", api.ApiSessionRequired(getPrivateChannelsForTeam)).Methods("GET")
	api.BaseRoutes.ChannelsForTeam.Handle("/ids", api.ApiSessionRequired(getPublicChannelsByIdsForTeam)).Methods("POST")
	api.BaseRoutes.ChannelsForTeam.Handle("/search", api.ApiSessionRequiredDisableWhenBusy(searchChannelsForTeam)).Methods("POST")
	api.BaseRoutes.ChannelsForTeam.Handle("/search_archived", api.ApiSessionRequiredDisableWhenBusy(searchArchivedChannelsForTeam)).Methods("POST")
	api.BaseRoutes.ChannelsForTeam.Handle("/autocomplete", api.ApiSessionRequired(autocompleteChannelsForTeam)).Methods("GET")
	api.BaseRoutes.ChannelsForTeam.Handle("/search_autocomplete", api.ApiSessionRequired(autocompleteChannelsForTeamForSearch)).Methods("GET")
	api.BaseRoutes.User.Handle("/teams/{team_id:[A-Za-z0-9]+}/channels", api.ApiSessionRequired(getChannelsForTeamForUser)).Methods("GET")

	api.BaseRoutes.ChannelCategories.Handle("", api.ApiSessionRequired(getCategoriesForTeamForUser)).Methods("GET")
	api.BaseRoutes.ChannelCategories.Handle("", api.ApiSessionRequired(createCategoryForTeamForUser)).Methods("POST")
	api.BaseRoutes.ChannelCategories.Handle("", api.ApiSessionRequired(updateCategoriesForTeamForUser)).Methods("PUT")
	api.BaseRoutes.ChannelCategories.Handle("/order", api.ApiSessionRequired(getCategoryOrderForTeamForUser)).Methods("GET")
	api.BaseRoutes.ChannelCategories.Handle("/order", api.ApiSessionRequired(updateCategoryOrderForTeamForUser)).Methods("PUT")
	api.BaseRoutes.ChannelCategories.Handle("/{category_id:[A-Za-z0-9_-]+}", api.ApiSessionRequired(getCategoryForTeamForUser)).Methods("GET")
	api.BaseRoutes.ChannelCategories.Handle("/{category_id:[A-Za-z0-9_-]+}", api.ApiSessionRequired(updateCategoryForTeamForUser)).Methods("PUT")
	api.BaseRoutes.ChannelCategories.Handle("/{category_id:[A-Za-z0-9_-]+}", api.ApiSessionRequired(deleteCategoryForTeamForUser)).Methods("DELETE")

	api.BaseRoutes.Channel.Handle("", api.ApiSessionRequired(getChannel)).Methods("GET")
	api.BaseRoutes.Channel.Handle("", api.ApiSessionRequired(updateChannel)).Methods("PUT")
	api.BaseRoutes.Channel.Handle("/patch", api.ApiSessionRequired(patchChannel)).Methods("PUT")
	api.BaseRoutes.Channel.Handle("/convert", api.ApiSessionRequired(convertChannelToPrivate)).Methods("POST")
	api.BaseRoutes.Channel.Handle("/privacy", api.ApiSessionRequired(updateChannelPrivacy)).Methods("PUT")
	api.BaseRoutes.Channel.Handle("/restore", api.ApiSessionRequired(restoreChannel)).Methods("POST")
	api.BaseRoutes.Channel.Handle("", api.ApiSessionRequired(deleteChannel)).Methods("DELETE")
	api.BaseRoutes.Channel.Handle("/stats", api.ApiSessionRequired(getChannelStats)).Methods("GET")
	api.BaseRoutes.Channel.Handle("/pinned", api.ApiSessionRequired(getPinnedPosts)).Methods("GET")
	api.BaseRoutes.Channel.Handle("/timezones", api.ApiSessionRequired(getChannelMembersTimezones)).Methods("GET")
	api.BaseRoutes.Channel.Handle("/members_minus_group_members", api.ApiSessionRequired(channelMembersMinusGroupMembers)).Methods("GET")
	api.BaseRoutes.Channel.Handle("/move", api.ApiSessionRequired(moveChannel)).Methods("POST")
	api.BaseRoutes.Channel.Handle("/member_counts_by_group", api.ApiSessionRequired(channelMemberCountsByGroup)).Methods("GET")

	api.BaseRoutes.ChannelForUser.Handle("/unread", api.ApiSessionRequired(getChannelUnread)).Methods("GET")

	api.BaseRoutes.ChannelByName.Handle("", api.ApiSessionRequired(getChannelByName)).Methods("GET")
	api.BaseRoutes.ChannelByNameForTeamName.Handle("", api.ApiSessionRequired(getChannelByNameForTeamName)).Methods("GET")

	api.BaseRoutes.ChannelMembers.Handle("", api.ApiSessionRequired(getChannelMembers)).Methods("GET")
	api.BaseRoutes.ChannelMembers.Handle("/ids", api.ApiSessionRequired(getChannelMembersByIds)).Methods("POST")
	api.BaseRoutes.ChannelMembers.Handle("", api.ApiSessionRequired(addChannelMember)).Methods("POST")
	api.BaseRoutes.ChannelMembersForUser.Handle("", api.ApiSessionRequired(getChannelMembersForUser)).Methods("GET")
	api.BaseRoutes.ChannelMember.Handle("", api.ApiSessionRequired(getChannelMember)).Methods("GET")
	api.BaseRoutes.ChannelMember.Handle("", api.ApiSessionRequired(removeChannelMember)).Methods("DELETE")
	api.BaseRoutes.ChannelMember.Handle("/roles", api.ApiSessionRequired(updateChannelMemberRoles)).Methods("PUT")
	api.BaseRoutes.ChannelMember.Handle("/schemeRoles", api.ApiSessionRequired(updateChannelMemberSchemeRoles)).Methods("PUT")
	api.BaseRoutes.ChannelMember.Handle("/notify_props", api.ApiSessionRequired(updateChannelMemberNotifyProps)).Methods("PUT")

	api.BaseRoutes.ChannelModerations.Handle("", api.ApiSessionRequired(getChannelModerations)).Methods("GET")
	api.BaseRoutes.ChannelModerations.Handle("/patch", api.ApiSessionRequired(patchChannelModerations)).Methods("PUT")
}

func createChannel(c *Context, w http.ResponseWriter, r *http.Request) {
	var channel *model.Channel
	err := json.NewDecoder(r.Body).Decode(&channel)
	if err != nil {
		c.SetInvalidParam("channel")
		return
	}

	auditRec := c.MakeAuditRecord("createChannel", audit.Fail)
	defer c.LogAuditRec(auditRec)
	auditRec.AddMeta("channel", channel)

	if channel.Type == model.ChannelTypeOpen && !c.App.SessionHasPermissionToTeam(*c.AppContext.Session(), channel.TeamId, model.PermissionCreatePublicChannel) {
		c.SetPermissionError(model.PermissionCreatePublicChannel)
		return
	}

	if channel.Type == model.ChannelTypePrivate && !c.App.SessionHasPermissionToTeam(*c.AppContext.Session(), channel.TeamId, model.PermissionCreatePrivateChannel) {
		c.SetPermissionError(model.PermissionCreatePrivateChannel)
		return
	}

	sc, appErr := c.App.CreateChannelWithUser(c.AppContext, channel, c.AppContext.Session().UserId)
	if appErr != nil {
		c.Err = appErr
		return
	}

	auditRec.Success()
	auditRec.AddMeta("channel", sc) // overwrite meta
	c.LogAudit("name=" + channel.Name)

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(sc); err != nil {
		mlog.Warn("Error while writing response", mlog.Err(err))
	}
}

func updateChannel(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireChannelId()
	if c.Err != nil {
		return
	}

	var channel *model.Channel
	err := json.NewDecoder(r.Body).Decode(&channel)
	if err != nil {
		c.SetInvalidParam("channel")
		return
	}

	// The channel being updated in the payload must be the same one as indicated in the URL.
	if channel.Id != c.Params.ChannelId {
		c.SetInvalidParam("channel_id")
		return
	}

	auditRec := c.MakeAuditRecord("updateChannel", audit.Fail)
	defer c.LogAuditRec(auditRec)

	originalOldChannel, appErr := c.App.GetChannel(channel.Id)
	if appErr != nil {
		c.Err = appErr
		return
	}
	oldChannel := originalOldChannel.DeepCopy()

	auditRec.AddMeta("channel", oldChannel)

	switch oldChannel.Type {
	case model.ChannelTypeOpen:
		if !c.App.SessionHasPermissionToChannel(*c.AppContext.Session(), c.Params.ChannelId, model.PermissionManagePublicChannelProperties) {
			c.SetPermissionError(model.PermissionManagePublicChannelProperties)
			return
		}

	case model.ChannelTypePrivate:
		if !c.App.SessionHasPermissionToChannel(*c.AppContext.Session(), c.Params.ChannelId, model.PermissionManagePrivateChannelProperties) {
			c.SetPermissionError(model.PermissionManagePrivateChannelProperties)
			return
		}

	case model.ChannelTypeGroup, model.ChannelTypeDirect:
		// Modifying the header is not linked to any specific permission for group/dm channels, so just check for membership.
		if _, errGet := c.App.GetChannelMember(context.Background(), channel.Id, c.AppContext.Session().UserId); errGet != nil {
			c.Err = model.NewAppError("updateChannel", "api.channel.patch_update_channel.forbidden.app_error", nil, "", http.StatusForbidden)
			return
		}

	default:
		c.Err = model.NewAppError("updateChannel", "api.channel.patch_update_channel.forbidden.app_error", nil, "", http.StatusForbidden)
		return
	}

	if oldChannel.DeleteAt > 0 {
		c.Err = model.NewAppError("updateChannel", "api.channel.update_channel.deleted.app_error", nil, "", http.StatusBadRequest)
		return
	}

	if channel.Type != "" && channel.Type != oldChannel.Type {
		c.Err = model.NewAppError("updateChannel", "api.channel.update_channel.typechange.app_error", nil, "", http.StatusBadRequest)
		return
	}

	if oldChannel.Name == model.DefaultChannelName {
		if channel.Name != "" && channel.Name != oldChannel.Name {
			c.Err = model.NewAppError("updateChannel", "api.channel.update_channel.tried.app_error", map[string]interface{}{"Channel": model.DefaultChannelName}, "", http.StatusBadRequest)
			return
		}
	}

	oldChannel.Header = channel.Header
	oldChannel.Purpose = channel.Purpose

	oldChannelDisplayName := oldChannel.DisplayName

	if channel.DisplayName != "" {
		oldChannel.DisplayName = channel.DisplayName
	}

	if channel.Name != "" {
		oldChannel.Name = channel.Name
		auditRec.AddMeta("new_channel_name", oldChannel.Name)
	}

	if channel.GroupConstrained != nil {
		oldChannel.GroupConstrained = channel.GroupConstrained
	}

	updatedChannel, appErr := c.App.UpdateChannel(oldChannel)
	if appErr != nil {
		c.Err = appErr
		return
	}
	auditRec.AddMeta("update", updatedChannel)

	if oldChannelDisplayName != channel.DisplayName {
		if err := c.App.PostUpdateChannelDisplayNameMessage(c.AppContext, c.AppContext.Session().UserId, channel, oldChannelDisplayName, channel.DisplayName); err != nil {
			mlog.Warn("Error while posting channel display name message", mlog.Err(err))
		}
	}

	auditRec.Success()
	c.LogAudit("name=" + channel.Name)

	if err := json.NewEncoder(w).Encode(oldChannel); err != nil {
		mlog.Warn("Error while writing response", mlog.Err(err))
	}
}

func convertChannelToPrivate(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireChannelId()
	if c.Err != nil {
		return
	}

	oldPublicChannel, err := c.App.GetChannel(c.Params.ChannelId)
	if err != nil {
		c.Err = err
		return
	}

	auditRec := c.MakeAuditRecord("convertChannelToPrivate", audit.Fail)
	defer c.LogAuditRec(auditRec)
	auditRec.AddMeta("channel", oldPublicChannel)

	if !c.App.SessionHasPermissionToChannel(*c.AppContext.Session(), c.Params.ChannelId, model.PermissionConvertPublicChannelToPrivate) {
		c.SetPermissionError(model.PermissionConvertPublicChannelToPrivate)
		return
	}

	if oldPublicChannel.Type == model.ChannelTypePrivate {
		c.Err = model.NewAppError("convertChannelToPrivate", "api.channel.convert_channel_to_private.private_channel_error", nil, "", http.StatusBadRequest)
		return
	}

	if oldPublicChannel.Name == model.DefaultChannelName {
		c.Err = model.NewAppError("convertChannelToPrivate", "api.channel.convert_channel_to_private.default_channel_error", nil, "", http.StatusBadRequest)
		return
	}

	user, err := c.App.GetUser(c.AppContext.Session().UserId)
	if err != nil {
		c.Err = err
		return
	}
	auditRec.AddMeta("user", user)

	oldPublicChannel.Type = model.ChannelTypePrivate

	rchannel, err := c.App.UpdateChannelPrivacy(c.AppContext, oldPublicChannel, user)
	if err != nil {
		c.Err = err
		return
	}

	auditRec.Success()
	c.LogAudit("name=" + rchannel.Name)

	if err := json.NewEncoder(w).Encode(rchannel); err != nil {
		mlog.Warn("Error while writing response", mlog.Err(err))
	}
}

func updateChannelPrivacy(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireChannelId()
	if c.Err != nil {
		return
	}

	props := model.StringInterfaceFromJson(r.Body)
	privacy, ok := props["privacy"].(string)
	if !ok || (model.ChannelType(privacy) != model.ChannelTypeOpen && model.ChannelType(privacy) != model.ChannelTypePrivate) {
		c.SetInvalidParam("privacy")
		return
	}

	channel, err := c.App.GetChannel(c.Params.ChannelId)
	if err != nil {
		c.Err = err
		return
	}

	auditRec := c.MakeAuditRecord("updateChannelPrivacy", audit.Fail)
	defer c.LogAuditRec(auditRec)
	auditRec.AddMeta("channel", channel)
	auditRec.AddMeta("new_type", privacy)

	if model.ChannelType(privacy) == model.ChannelTypeOpen && !c.App.SessionHasPermissionToChannel(*c.AppContext.Session(), c.Params.ChannelId, model.PermissionConvertPrivateChannelToPublic) {
		c.SetPermissionError(model.PermissionConvertPrivateChannelToPublic)
		return
	}

	if model.ChannelType(privacy) == model.ChannelTypePrivate && !c.App.SessionHasPermissionToChannel(*c.AppContext.Session(), c.Params.ChannelId, model.PermissionConvertPublicChannelToPrivate) {
		c.SetPermissionError(model.PermissionConvertPublicChannelToPrivate)
		return
	}

	if channel.Name == model.DefaultChannelName && model.ChannelType(privacy) == model.ChannelTypePrivate {
		c.Err = model.NewAppError("updateChannelPrivacy", "api.channel.update_channel_privacy.default_channel_error", nil, "", http.StatusBadRequest)
		return
	}

	user, err := c.App.GetUser(c.AppContext.Session().UserId)
	if err != nil {
		c.Err = err
		return
	}
	auditRec.AddMeta("user", user)

	channel.Type = model.ChannelType(privacy)

	updatedChannel, err := c.App.UpdateChannelPrivacy(c.AppContext, channel, user)
	if err != nil {
		c.Err = err
		return
	}

	auditRec.Success()
	c.LogAudit("name=" + updatedChannel.Name)

	if err := json.NewEncoder(w).Encode(updatedChannel); err != nil {
		mlog.Warn("Error while writing response", mlog.Err(err))
	}
}

func patchChannel(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireChannelId()
	if c.Err != nil {
		return
	}
	var patch *model.ChannelPatch
	err := json.NewDecoder(r.Body).Decode(&patch)
	if err != nil {
		c.SetInvalidParam("channel")
		return
	}

	originalOldChannel, appErr := c.App.GetChannel(c.Params.ChannelId)
	if appErr != nil {
		c.Err = appErr
		return
	}
	oldChannel := originalOldChannel.DeepCopy()

	auditRec := c.MakeAuditRecord("patchChannel", audit.Fail)
	defer c.LogAuditRec(auditRec)
	auditRec.AddMeta("channel", oldChannel)

	switch oldChannel.Type {
	case model.ChannelTypeOpen:
		if !c.App.SessionHasPermissionToChannel(*c.AppContext.Session(), c.Params.ChannelId, model.PermissionManagePublicChannelProperties) {
			c.SetPermissionError(model.PermissionManagePublicChannelProperties)
			return
		}

	case model.ChannelTypePrivate:
		if !c.App.SessionHasPermissionToChannel(*c.AppContext.Session(), c.Params.ChannelId, model.PermissionManagePrivateChannelProperties) {
			c.SetPermissionError(model.PermissionManagePrivateChannelProperties)
			return
		}

	case model.ChannelTypeGroup, model.ChannelTypeDirect:
		// Modifying the header is not linked to any specific permission for group/dm channels, so just check for membership.
		if _, appErr = c.App.GetChannelMember(context.Background(), c.Params.ChannelId, c.AppContext.Session().UserId); appErr != nil {
			c.Err = model.NewAppError("patchChannel", "api.channel.patch_update_channel.forbidden.app_error", nil, "", http.StatusForbidden)
			return
		}

	default:
		c.Err = model.NewAppError("patchChannel", "api.channel.patch_update_channel.forbidden.app_error", nil, "", http.StatusForbidden)
		return
	}

	rchannel, appErr := c.App.PatchChannel(c.AppContext, oldChannel, patch, c.AppContext.Session().UserId)
	if appErr != nil {
		c.Err = appErr
		return
	}

	appErr = c.App.FillInChannelProps(rchannel)
	if appErr != nil {
		c.Err = appErr
		return
	}

	auditRec.Success()
	c.LogAudit("")
	auditRec.AddMeta("patch", rchannel)

	if err := json.NewEncoder(w).Encode(rchannel); err != nil {
		mlog.Warn("Error while writing response", mlog.Err(err))
	}
}

func restoreChannel(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireChannelId()
	if c.Err != nil {
		return
	}

	channel, err := c.App.GetChannel(c.Params.ChannelId)
	if err != nil {
		c.Err = err
		return
	}
	teamId := channel.TeamId

	auditRec := c.MakeAuditRecord("restoreChannel", audit.Fail)
	defer c.LogAuditRec(auditRec)
	auditRec.AddMeta("channel", channel)

	if !c.App.SessionHasPermissionToTeam(*c.AppContext.Session(), teamId, model.PermissionManageTeam) {
		c.SetPermissionError(model.PermissionManageTeam)
		return
	}

	channel, err = c.App.RestoreChannel(c.AppContext, channel, c.AppContext.Session().UserId)
	if err != nil {
		c.Err = err
		return
	}

	auditRec.Success()
	c.LogAudit("name=" + channel.Name)

	if err := json.NewEncoder(w).Encode(channel); err != nil {
		mlog.Warn("Error while writing response", mlog.Err(err))
	}
}

func createDirectChannel(c *Context, w http.ResponseWriter, r *http.Request) {
	userIds := model.ArrayFromJson(r.Body)
	allowed := false

	if len(userIds) != 2 {
		c.SetInvalidParam("user_ids")
		return
	}

	for _, id := range userIds {
		if !model.IsValidId(id) {
			c.SetInvalidParam("user_id")
			return
		}
		if id == c.AppContext.Session().UserId {
			allowed = true
		}
	}

	auditRec := c.MakeAuditRecord("createDirectChannel", audit.Fail)
	defer c.LogAuditRec(auditRec)

	if !c.App.SessionHasPermissionTo(*c.AppContext.Session(), model.PermissionCreateDirectChannel) {
		c.SetPermissionError(model.PermissionCreateDirectChannel)
		return
	}

	if !allowed && !c.App.SessionHasPermissionTo(*c.AppContext.Session(), model.PermissionManageSystem) {
		c.SetPermissionError(model.PermissionManageSystem)
		return
	}

	otherUserId := userIds[0]
	if c.AppContext.Session().UserId == otherUserId {
		otherUserId = userIds[1]
	}

	auditRec.AddMeta("other_user_id", otherUserId)

	canSee, err := c.App.UserCanSeeOtherUser(c.AppContext.Session().UserId, otherUserId)
	if err != nil {
		c.Err = err
		return
	}

	if !canSee {
		c.SetPermissionError(model.PermissionViewMembers)
		return
	}

	sc, err := c.App.GetOrCreateDirectChannel(c.AppContext, userIds[0], userIds[1])
	if err != nil {
		c.Err = err
		return
	}

	auditRec.Success()
	auditRec.AddMeta("channel", sc)

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(sc); err != nil {
		mlog.Warn("Error while writing response", mlog.Err(err))
	}
}

func searchGroupChannels(c *Context, w http.ResponseWriter, r *http.Request) {
	var props *model.ChannelSearch
	err := json.NewDecoder(r.Body).Decode(&props)
	if err != nil {
		c.SetInvalidParam("channel_search")
		return
	}

	groupChannels, appErr := c.App.SearchGroupChannels(c.AppContext.Session().UserId, props.Term)
	if appErr != nil {
		c.Err = appErr
		return
	}

	if err := json.NewEncoder(w).Encode(groupChannels); err != nil {
		mlog.Warn("Error while writing response", mlog.Err(err))
	}
}

func createGroupChannel(c *Context, w http.ResponseWriter, r *http.Request) {
	userIds := model.ArrayFromJson(r.Body)

	if len(userIds) == 0 {
		c.SetInvalidParam("user_ids")
		return
	}

	found := false
	for _, id := range userIds {
		if !model.IsValidId(id) {
			c.SetInvalidParam("user_id")
			return
		}
		if id == c.AppContext.Session().UserId {
			found = true
		}
	}

	if !found {
		userIds = append(userIds, c.AppContext.Session().UserId)
	}

	auditRec := c.MakeAuditRecord("createGroupChannel", audit.Fail)
	defer c.LogAuditRec(auditRec)

	if !c.App.SessionHasPermissionTo(*c.AppContext.Session(), model.PermissionCreateGroupChannel) {
		c.SetPermissionError(model.PermissionCreateGroupChannel)
		return
	}

	canSeeAll := true
	for _, id := range userIds {
		if c.AppContext.Session().UserId != id {
			canSee, err := c.App.UserCanSeeOtherUser(c.AppContext.Session().UserId, id)
			if err != nil {
				c.Err = err
				return
			}
			if !canSee {
				canSeeAll = false
			}
		}
	}

	if !canSeeAll {
		c.SetPermissionError(model.PermissionViewMembers)
		return
	}

	groupChannel, err := c.App.CreateGroupChannel(userIds, c.AppContext.Session().UserId)
	if err != nil {
		c.Err = err
		return
	}

	auditRec.Success()
	auditRec.AddMeta("channel", groupChannel)

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(groupChannel); err != nil {
		mlog.Warn("Error while writing response", mlog.Err(err))
	}
}

func getChannel(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireChannelId()
	if c.Err != nil {
		return
	}

	channel, err := c.App.GetChannel(c.Params.ChannelId)
	if err != nil {
		c.Err = err
		return
	}

	if channel.Type == model.ChannelTypeOpen {
		if !c.App.SessionHasPermissionToTeam(*c.AppContext.Session(), channel.TeamId, model.PermissionReadPublicChannel) && !c.App.SessionHasPermissionToChannel(*c.AppContext.Session(), c.Params.ChannelId, model.PermissionReadChannel) {
			c.SetPermissionError(model.PermissionReadPublicChannel)
			return
		}
	} else {
		if !c.App.SessionHasPermissionToChannel(*c.AppContext.Session(), c.Params.ChannelId, model.PermissionReadChannel) {
			c.SetPermissionError(model.PermissionReadChannel)
			return
		}
	}

	err = c.App.FillInChannelProps(channel)
	if err != nil {
		c.Err = err
		return
	}

	if err := json.NewEncoder(w).Encode(channel); err != nil {
		mlog.Warn("Error while writing response", mlog.Err(err))
	}
}

func getChannelUnread(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireChannelId().RequireUserId()
	if c.Err != nil {
		return
	}

	if !c.App.SessionHasPermissionToUser(*c.AppContext.Session(), c.Params.UserId) {
		c.SetPermissionError(model.PermissionEditOtherUsers)
		return
	}

	if !c.App.SessionHasPermissionToChannel(*c.AppContext.Session(), c.Params.ChannelId, model.PermissionReadChannel) {
		c.SetPermissionError(model.PermissionReadChannel)
		return
	}

	channelUnread, err := c.App.GetChannelUnread(c.Params.ChannelId, c.Params.UserId)
	if err != nil {
		c.Err = err
		return
	}

	if err := json.NewEncoder(w).Encode(channelUnread); err != nil {
		mlog.Warn("Error while writing response", mlog.Err(err))
	}
}

func getChannelStats(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireChannelId()
	if c.Err != nil {
		return
	}

	if !c.App.SessionHasPermissionToChannel(*c.AppContext.Session(), c.Params.ChannelId, model.PermissionReadChannel) {
		c.SetPermissionError(model.PermissionReadChannel)
		return
	}

	memberCount, err := c.App.GetChannelMemberCount(c.Params.ChannelId)
	if err != nil {
		c.Err = err
		return
	}

	guestCount, err := c.App.GetChannelGuestCount(c.Params.ChannelId)
	if err != nil {
		c.Err = err
		return
	}

	pinnedPostCount, err := c.App.GetChannelPinnedPostCount(c.Params.ChannelId)
	if err != nil {
		c.Err = err
		return
	}

	stats := model.ChannelStats{
		ChannelId:       c.Params.ChannelId,
		MemberCount:     memberCount,
		GuestCount:      guestCount,
		PinnedPostCount: pinnedPostCount,
	}
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		mlog.Warn("Error while writing response", mlog.Err(err))
	}
}

func getPinnedPosts(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireChannelId()
	if c.Err != nil {
		return
	}

	if !c.App.SessionHasPermissionToChannel(*c.AppContext.Session(), c.Params.ChannelId, model.PermissionReadChannel) {
		c.SetPermissionError(model.PermissionReadChannel)
		return
	}

	posts, err := c.App.GetPinnedPosts(c.Params.ChannelId)
	if err != nil {
		c.Err = err
		return
	}

	if c.HandleEtag(posts.Etag(), "Get Pinned Posts", w, r) {
		return
	}

	clientPostList := c.App.PreparePostListForClient(posts)

	w.Header().Set(model.HeaderEtagServer, clientPostList.Etag())
	if err := json.NewEncoder(w).Encode(clientPostList); err != nil {
		mlog.Warn("Error while writing response", mlog.Err(err))
	}
}

func getAllChannels(c *Context, w http.ResponseWriter, r *http.Request) {
	permissions := []*model.Permission{
		model.PermissionSysconsoleReadUserManagementGroups,
		model.PermissionSysconsoleReadUserManagementChannels,
	}
	if !c.App.SessionHasPermissionToAny(*c.AppContext.Session(), permissions) {
		c.SetPermissionError(permissions...)
		return
	}
	// Only system managers may use the ExcludePolicyConstrained parameter
	if c.Params.ExcludePolicyConstrained && !c.App.SessionHasPermissionTo(*c.AppContext.Session(), model.PermissionSysconsoleReadComplianceDataRetentionPolicy) {
		c.SetPermissionError(model.PermissionSysconsoleReadComplianceDataRetentionPolicy)
		return
	}

	opts := model.ChannelSearchOpts{
		NotAssociatedToGroup:     c.Params.NotAssociatedToGroup,
		ExcludeDefaultChannels:   c.Params.ExcludeDefaultChannels,
		IncludeDeleted:           c.Params.IncludeDeleted,
		ExcludePolicyConstrained: c.Params.ExcludePolicyConstrained,
	}
	if c.App.SessionHasPermissionTo(*c.AppContext.Session(), model.PermissionSysconsoleReadComplianceDataRetentionPolicy) {
		opts.IncludePolicyID = true
	}

	channels, err := c.App.GetAllChannels(c.Params.Page, c.Params.PerPage, opts)
	if err != nil {
		c.Err = err
		return
	}

	if c.Params.IncludeTotalCount {
		totalCount, err := c.App.GetAllChannelsCount(opts)
		if err != nil {
			c.Err = err
			return
		}
		cwc := &model.ChannelsWithCount{
			Channels:   channels,
			TotalCount: totalCount,
		}
		if err := json.NewEncoder(w).Encode(cwc); err != nil {
			mlog.Warn("Error while writing response", mlog.Err(err))
		}
		return
	}

	if err := json.NewEncoder(w).Encode(channels); err != nil {
		mlog.Warn("Error while writing response", mlog.Err(err))
	}
}

func getPublicChannelsForTeam(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireTeamId()
	if c.Err != nil {
		return
	}

	if !c.App.SessionHasPermissionToTeam(*c.AppContext.Session(), c.Params.TeamId, model.PermissionListTeamChannels) {
		c.SetPermissionError(model.PermissionListTeamChannels)
		return
	}

	channels, err := c.App.GetPublicChannelsForTeam(c.Params.TeamId, c.Params.Page*c.Params.PerPage, c.Params.PerPage)
	if err != nil {
		c.Err = err
		return
	}

	err = c.App.FillInChannelsProps(channels)
	if err != nil {
		c.Err = err
		return
	}

	if err := json.NewEncoder(w).Encode(channels); err != nil {
		mlog.Warn("Error while writing response", mlog.Err(err))
	}
}

func getDeletedChannelsForTeam(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireTeamId()
	if c.Err != nil {
		return
	}

	channels, err := c.App.GetDeletedChannels(c.Params.TeamId, c.Params.Page*c.Params.PerPage, c.Params.PerPage, c.AppContext.Session().UserId)
	if err != nil {
		c.Err = err
		return
	}

	err = c.App.FillInChannelsProps(channels)
	if err != nil {
		c.Err = err
		return
	}

	if err := json.NewEncoder(w).Encode(channels); err != nil {
		mlog.Warn("Error while writing response", mlog.Err(err))
	}
}

func getPrivateChannelsForTeam(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireTeamId()
	if c.Err != nil {
		return
	}

	if !c.App.SessionHasPermissionTo(*c.AppContext.Session(), model.PermissionManageSystem) {
		c.SetPermissionError(model.PermissionManageSystem)
		return
	}

	channels, err := c.App.GetPrivateChannelsForTeam(c.Params.TeamId, c.Params.Page*c.Params.PerPage, c.Params.PerPage)
	if err != nil {
		c.Err = err
		return
	}

	err = c.App.FillInChannelsProps(channels)
	if err != nil {
		c.Err = err
		return
	}

	if err := json.NewEncoder(w).Encode(channels); err != nil {
		mlog.Warn("Error while writing response", mlog.Err(err))
	}
}

func getPublicChannelsByIdsForTeam(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireTeamId()
	if c.Err != nil {
		return
	}

	channelIds := model.ArrayFromJson(r.Body)
	if len(channelIds) == 0 {
		c.SetInvalidParam("channel_ids")
		return
	}

	for _, cid := range channelIds {
		if !model.IsValidId(cid) {
			c.SetInvalidParam("channel_id")
			return
		}
	}

	if !c.App.SessionHasPermissionToTeam(*c.AppContext.Session(), c.Params.TeamId, model.PermissionViewTeam) {
		c.SetPermissionError(model.PermissionViewTeam)
		return
	}

	channels, err := c.App.GetPublicChannelsByIdsForTeam(c.Params.TeamId, channelIds)
	if err != nil {
		c.Err = err
		return
	}

	err = c.App.FillInChannelsProps(channels)
	if err != nil {
		c.Err = err
		return
	}

	if err := json.NewEncoder(w).Encode(channels); err != nil {
		mlog.Warn("Error while writing response", mlog.Err(err))
	}
}

func getChannelsForTeamForUser(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireUserId().RequireTeamId()
	if c.Err != nil {
		return
	}

	if !c.App.SessionHasPermissionToUser(*c.AppContext.Session(), c.Params.UserId) {
		c.SetPermissionError(model.PermissionEditOtherUsers)
		return
	}

	if !c.App.SessionHasPermissionToTeam(*c.AppContext.Session(), c.Params.TeamId, model.PermissionViewTeam) {
		c.SetPermissionError(model.PermissionViewTeam)
		return
	}

	query := r.URL.Query()
	lastDeleteAt, nErr := strconv.Atoi(query.Get("last_delete_at"))
	if nErr != nil {
		lastDeleteAt = 0
	}
	if lastDeleteAt < 0 {
		c.SetInvalidUrlParam("last_delete_at")
		return
	}

	channels, err := c.App.GetChannelsForUser(c.Params.TeamId, c.Params.UserId, c.Params.IncludeDeleted, lastDeleteAt)
	if err != nil {
		c.Err = err
		return
	}

	if c.HandleEtag(channels.Etag(), "Get Channels", w, r) {
		return
	}

	err = c.App.FillInChannelsProps(channels)
	if err != nil {
		c.Err = err
		return
	}

	w.Header().Set(model.HeaderEtagServer, channels.Etag())
	if err := json.NewEncoder(w).Encode(channels); err != nil {
		mlog.Warn("Error while writing response", mlog.Err(err))
	}
}

func autocompleteChannelsForTeam(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireTeamId()
	if c.Err != nil {
		return
	}

	if !c.App.SessionHasPermissionToTeam(*c.AppContext.Session(), c.Params.TeamId, model.PermissionListTeamChannels) {
		c.SetPermissionError(model.PermissionListTeamChannels)
		return
	}

	name := r.URL.Query().Get("name")

	channels, err := c.App.AutocompleteChannels(c.Params.TeamId, name)
	if err != nil {
		c.Err = err
		return
	}

	// Don't fill in channels props, since unused by client and potentially expensive.

	if err := json.NewEncoder(w).Encode(channels); err != nil {
		mlog.Warn("Error while writing response", mlog.Err(err))
	}
}

func autocompleteChannelsForTeamForSearch(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireTeamId()
	if c.Err != nil {
		return
	}

	name := r.URL.Query().Get("name")

	channels, err := c.App.AutocompleteChannelsForSearch(c.Params.TeamId, c.AppContext.Session().UserId, name)
	if err != nil {
		c.Err = err
		return
	}

	if err := json.NewEncoder(w).Encode(channels); err != nil {
		mlog.Warn("Error while writing response", mlog.Err(err))
	}
}

func searchChannelsForTeam(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireTeamId()
	if c.Err != nil {
		return
	}

	var props *model.ChannelSearch
	err := json.NewDecoder(r.Body).Decode(&props)
	if err != nil {
		c.SetInvalidParam("channel_search")
		return
	}

	var channels *model.ChannelList
	var appErr *model.AppError
	if c.App.SessionHasPermissionToTeam(*c.AppContext.Session(), c.Params.TeamId, model.PermissionListTeamChannels) {
		channels, appErr = c.App.SearchChannels(c.Params.TeamId, props.Term)
	} else {
		// If the user is not a team member, return a 404
		if _, appErr = c.App.GetTeamMember(c.Params.TeamId, c.AppContext.Session().UserId); appErr != nil {
			c.Err = appErr
			return
		}

		channels, appErr = c.App.SearchChannelsForUser(c.AppContext.Session().UserId, c.Params.TeamId, props.Term)
	}

	if appErr != nil {
		c.Err = appErr
		return
	}

	// Don't fill in channels props, since unused by client and potentially expensive.

	if err := json.NewEncoder(w).Encode(channels); err != nil {
		mlog.Warn("Error while writing response", mlog.Err(err))
	}
}

func searchArchivedChannelsForTeam(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireTeamId()
	if c.Err != nil {
		return
	}

	var props *model.ChannelSearch
	err := json.NewDecoder(r.Body).Decode(&props)
	if err != nil {
		c.SetInvalidParam("channel_search")
		return
	}

	var channels *model.ChannelList
	var appErr *model.AppError
	if c.App.SessionHasPermissionToTeam(*c.AppContext.Session(), c.Params.TeamId, model.PermissionListTeamChannels) {
		channels, appErr = c.App.SearchArchivedChannels(c.Params.TeamId, props.Term, c.AppContext.Session().UserId)
	} else {
		// If the user is not a team member, return a 404
		if _, appErr = c.App.GetTeamMember(c.Params.TeamId, c.AppContext.Session().UserId); appErr != nil {
			c.Err = appErr
			return
		}

		channels, appErr = c.App.SearchArchivedChannels(c.Params.TeamId, props.Term, c.AppContext.Session().UserId)
	}

	if appErr != nil {
		c.Err = appErr
		return
	}

	// Don't fill in channels props, since unused by client and potentially expensive.

	if err := json.NewEncoder(w).Encode(channels); err != nil {
		mlog.Warn("Error while writing response", mlog.Err(err))
	}
}

func searchAllChannels(c *Context, w http.ResponseWriter, r *http.Request) {
	var props *model.ChannelSearch
	err := json.NewDecoder(r.Body).Decode(&props)
	if err != nil {
		c.SetInvalidParam("channel_search")
		return
	}
	// Only system managers may use the ExcludePolicyConstrained field
	if props.ExcludePolicyConstrained && !c.App.SessionHasPermissionTo(*c.AppContext.Session(), model.PermissionSysconsoleReadComplianceDataRetentionPolicy) {
		c.SetPermissionError(model.PermissionSysconsoleReadComplianceDataRetentionPolicy)
		return
	}

	if !c.App.SessionHasPermissionTo(*c.AppContext.Session(), model.PermissionSysconsoleReadUserManagementChannels) {
		c.SetPermissionError(model.PermissionSysconsoleReadUserManagementChannels)
		return
	}
	includeDeleted, _ := strconv.ParseBool(r.URL.Query().Get("include_deleted"))
	includeDeleted = includeDeleted || props.IncludeDeleted

	opts := model.ChannelSearchOpts{
		NotAssociatedToGroup:     props.NotAssociatedToGroup,
		ExcludeDefaultChannels:   props.ExcludeDefaultChannels,
		TeamIds:                  props.TeamIds,
		GroupConstrained:         props.GroupConstrained,
		ExcludeGroupConstrained:  props.ExcludeGroupConstrained,
		ExcludePolicyConstrained: props.ExcludePolicyConstrained,
		Public:                   props.Public,
		Private:                  props.Private,
		IncludeDeleted:           includeDeleted,
		Deleted:                  props.Deleted,
		Page:                     props.Page,
		PerPage:                  props.PerPage,
	}
	if c.App.SessionHasPermissionTo(*c.AppContext.Session(), model.PermissionSysconsoleReadComplianceDataRetentionPolicy) {
		opts.IncludePolicyID = true
	}

	channels, totalCount, appErr := c.App.SearchAllChannels(props.Term, opts)
	if appErr != nil {
		c.Err = appErr
		return
	}

	// Don't fill in channels props, since unused by client and potentially expensive.
	if props.Page != nil && props.PerPage != nil {
		data := model.ChannelsWithCount{Channels: channels, TotalCount: totalCount}

		if err := json.NewEncoder(w).Encode(data); err != nil {
			mlog.Warn("Error while writing response", mlog.Err(err))
		}
		return
	}

	if err := json.NewEncoder(w).Encode(channels); err != nil {
		mlog.Warn("Error while writing response", mlog.Err(err))
	}
}

func deleteChannel(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireChannelId()
	if c.Err != nil {
		return
	}

	channel, err := c.App.GetChannel(c.Params.ChannelId)
	if err != nil {
		c.Err = err
		return
	}

	auditRec := c.MakeAuditRecord("deleteChannel", audit.Fail)
	defer c.LogAuditRec(auditRec)
	auditRec.AddMeta("channeld", channel)

	if channel.Type == model.ChannelTypeDirect || channel.Type == model.ChannelTypeGroup {
		c.Err = model.NewAppError("deleteChannel", "api.channel.delete_channel.type.invalid", nil, "", http.StatusBadRequest)
		return
	}

	if channel.Type == model.ChannelTypeOpen && !c.App.SessionHasPermissionToChannel(*c.AppContext.Session(), channel.Id, model.PermissionDeletePublicChannel) {
		c.SetPermissionError(model.PermissionDeletePublicChannel)
		return
	}

	if channel.Type == model.ChannelTypePrivate && !c.App.SessionHasPermissionToChannel(*c.AppContext.Session(), channel.Id, model.PermissionDeletePrivateChannel) {
		c.SetPermissionError(model.PermissionDeletePrivateChannel)
		return
	}

	if c.Params.Permanent {
		if *c.App.Config().ServiceSettings.EnableAPIChannelDeletion {
			err = c.App.PermanentDeleteChannel(channel)
		} else {
			err = model.NewAppError("deleteChannel", "api.user.delete_channel.not_enabled.app_error", nil, "channelId="+c.Params.ChannelId, http.StatusUnauthorized)
		}
	} else {
		err = c.App.DeleteChannel(c.AppContext, channel, c.AppContext.Session().UserId)
	}
	if err != nil {
		c.Err = err
		return
	}

	auditRec.Success()
	c.LogAudit("name=" + channel.Name)

	ReturnStatusOK(w)
}

func getChannelByName(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireTeamId().RequireChannelName()
	if c.Err != nil {
		return
	}

	includeDeleted, _ := strconv.ParseBool(r.URL.Query().Get("include_deleted"))
	channel, appErr := c.App.GetChannelByName(c.Params.ChannelName, c.Params.TeamId, includeDeleted)
	if appErr != nil {
		c.Err = appErr
		return
	}

	if channel.Type == model.ChannelTypeOpen {
		if !c.App.SessionHasPermissionToTeam(*c.AppContext.Session(), channel.TeamId, model.PermissionReadPublicChannel) && !c.App.SessionHasPermissionToChannel(*c.AppContext.Session(), channel.Id, model.PermissionReadChannel) {
			c.SetPermissionError(model.PermissionReadPublicChannel)
			return
		}
	} else {
		if !c.App.SessionHasPermissionToChannel(*c.AppContext.Session(), channel.Id, model.PermissionReadChannel) {
			c.Err = model.NewAppError("getChannelByName", "app.channel.get_by_name.missing.app_error", nil, "teamId="+channel.TeamId+", "+"name="+channel.Name+"", http.StatusNotFound)
			return
		}
	}

	appErr = c.App.FillInChannelProps(channel)
	if appErr != nil {
		c.Err = appErr
		return
	}

	if err := json.NewEncoder(w).Encode(channel); err != nil {
		mlog.Warn("Error while writing response", mlog.Err(err))
	}
}

func getChannelByNameForTeamName(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireTeamName().RequireChannelName()
	if c.Err != nil {
		return
	}

	includeDeleted, _ := strconv.ParseBool(r.URL.Query().Get("include_deleted"))
	channel, appErr := c.App.GetChannelByNameForTeamName(c.Params.ChannelName, c.Params.TeamName, includeDeleted)
	if appErr != nil {
		c.Err = appErr
		return
	}

	teamOk := c.App.SessionHasPermissionToTeam(*c.AppContext.Session(), channel.TeamId, model.PermissionReadPublicChannel)
	channelOk := c.App.SessionHasPermissionToChannel(*c.AppContext.Session(), channel.Id, model.PermissionReadChannel)

	if channel.Type == model.ChannelTypeOpen {
		if !teamOk && !channelOk {
			c.SetPermissionError(model.PermissionReadPublicChannel)
			return
		}
	} else if !channelOk {
		c.Err = model.NewAppError("getChannelByNameForTeamName", "app.channel.get_by_name.missing.app_error", nil, "teamId="+channel.TeamId+", "+"name="+channel.Name+"", http.StatusNotFound)
		return
	}

	appErr = c.App.FillInChannelProps(channel)
	if appErr != nil {
		c.Err = appErr
		return
	}

	if err := json.NewEncoder(w).Encode(channel); err != nil {
		mlog.Warn("Error while writing response", mlog.Err(err))
	}
}

func getChannelMembers(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireChannelId()
	if c.Err != nil {
		return
	}

	if !c.App.SessionHasPermissionToChannel(*c.AppContext.Session(), c.Params.ChannelId, model.PermissionReadChannel) {
		c.SetPermissionError(model.PermissionReadChannel)
		return
	}

	members, err := c.App.GetChannelMembersPage(c.Params.ChannelId, c.Params.Page, c.Params.PerPage)
	if err != nil {
		c.Err = err
		return
	}

	if err := json.NewEncoder(w).Encode(members); err != nil {
		mlog.Warn("Error while writing response", mlog.Err(err))
	}
}

func getChannelMembersTimezones(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireChannelId()
	if c.Err != nil {
		return
	}

	if !c.App.SessionHasPermissionToChannel(*c.AppContext.Session(), c.Params.ChannelId, model.PermissionReadChannel) {
		c.SetPermissionError(model.PermissionReadChannel)
		return
	}

	membersTimezones, err := c.App.GetChannelMembersTimezones(c.Params.ChannelId)
	if err != nil {
		c.Err = err
		return
	}

	w.Write([]byte(model.ArrayToJson(membersTimezones)))
}

func getChannelMembersByIds(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireChannelId()
	if c.Err != nil {
		return
	}

	userIds := model.ArrayFromJson(r.Body)
	if len(userIds) == 0 {
		c.SetInvalidParam("user_ids")
		return
	}

	if !c.App.SessionHasPermissionToChannel(*c.AppContext.Session(), c.Params.ChannelId, model.PermissionReadChannel) {
		c.SetPermissionError(model.PermissionReadChannel)
		return
	}

	members, err := c.App.GetChannelMembersByIds(c.Params.ChannelId, userIds)
	if err != nil {
		c.Err = err
		return
	}

	if err := json.NewEncoder(w).Encode(members); err != nil {
		mlog.Warn("Error while writing response", mlog.Err(err))
	}
}

func getChannelMember(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireChannelId().RequireUserId()
	if c.Err != nil {
		return
	}

	if !c.App.SessionHasPermissionToChannel(*c.AppContext.Session(), c.Params.ChannelId, model.PermissionReadChannel) {
		c.SetPermissionError(model.PermissionReadChannel)
		return
	}

	member, err := c.App.GetChannelMember(app.WithMaster(context.Background()), c.Params.ChannelId, c.Params.UserId)
	if err != nil {
		c.Err = err
		return
	}

	if err := json.NewEncoder(w).Encode(member); err != nil {
		mlog.Warn("Error while writing response", mlog.Err(err))
	}
}

func getChannelMembersForUser(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireUserId().RequireTeamId()
	if c.Err != nil {
		return
	}

	if !c.App.SessionHasPermissionToTeam(*c.AppContext.Session(), c.Params.TeamId, model.PermissionViewTeam) {
		c.SetPermissionError(model.PermissionViewTeam)
		return
	}

	if c.AppContext.Session().UserId != c.Params.UserId && !c.App.SessionHasPermissionToTeam(*c.AppContext.Session(), c.Params.TeamId, model.PermissionManageSystem) {
		c.SetPermissionError(model.PermissionManageSystem)
		return
	}

	members, err := c.App.GetChannelMembersForUser(c.Params.TeamId, c.Params.UserId)
	if err != nil {
		c.Err = err
		return
	}

	if err := json.NewEncoder(w).Encode(members); err != nil {
		mlog.Warn("Error while writing response", mlog.Err(err))
	}
}

func viewChannel(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireUserId()
	if c.Err != nil {
		return
	}

	if !c.App.SessionHasPermissionToUser(*c.AppContext.Session(), c.Params.UserId) {
		c.SetPermissionError(model.PermissionEditOtherUsers)
		return
	}

	view := model.ChannelViewFromJson(r.Body)
	if view == nil {
		c.SetInvalidParam("channel_view")
		return
	}

	// Validate view struct
	// Check IDs are valid or blank. Blank IDs are used to denote focus loss or initial channel view.
	if view.ChannelId != "" && !model.IsValidId(view.ChannelId) {
		c.SetInvalidParam("channel_view.channel_id")
		return
	}
	if view.PrevChannelId != "" && !model.IsValidId(view.PrevChannelId) {
		c.SetInvalidParam("channel_view.prev_channel_id")
		return
	}

	times, err := c.App.ViewChannel(view, c.Params.UserId, c.AppContext.Session().Id, view.CollapsedThreadsSupported)
	if err != nil {
		c.Err = err
		return
	}

	c.App.UpdateLastActivityAtIfNeeded(*c.AppContext.Session())
	c.ExtendSessionExpiryIfNeeded(w, r)

	// Returning {"status": "OK", ...} for backwards compatibility
	resp := &model.ChannelViewResponse{
		Status:            "OK",
		LastViewedAtTimes: times,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		mlog.Warn("Error while writing response", mlog.Err(err))
	}
}

func updateChannelMemberRoles(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireChannelId().RequireUserId()
	if c.Err != nil {
		return
	}

	props := model.MapFromJson(r.Body)

	newRoles := props["roles"]
	if !(model.IsValidUserRoles(newRoles)) {
		c.SetInvalidParam("roles")
		return
	}

	auditRec := c.MakeAuditRecord("updateChannelMemberRoles", audit.Fail)
	defer c.LogAuditRec(auditRec)
	auditRec.AddMeta("channel_id", c.Params.ChannelId)
	auditRec.AddMeta("roles", newRoles)

	if !c.App.SessionHasPermissionToChannel(*c.AppContext.Session(), c.Params.ChannelId, model.PermissionManageChannelRoles) {
		c.SetPermissionError(model.PermissionManageChannelRoles)
		return
	}

	if _, err := c.App.UpdateChannelMemberRoles(c.Params.ChannelId, c.Params.UserId, newRoles); err != nil {
		c.Err = err
		return
	}

	auditRec.Success()

	ReturnStatusOK(w)
}

func updateChannelMemberSchemeRoles(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireChannelId().RequireUserId()
	if c.Err != nil {
		return
	}

	schemeRoles := model.SchemeRolesFromJson(r.Body)
	if schemeRoles == nil {
		c.SetInvalidParam("scheme_roles")
		return
	}

	auditRec := c.MakeAuditRecord("updateChannelMemberSchemeRoles", audit.Fail)
	defer c.LogAuditRec(auditRec)
	auditRec.AddMeta("channel_id", c.Params.ChannelId)
	auditRec.AddMeta("roles", schemeRoles)

	if !c.App.SessionHasPermissionToChannel(*c.AppContext.Session(), c.Params.ChannelId, model.PermissionManageChannelRoles) {
		c.SetPermissionError(model.PermissionManageChannelRoles)
		return
	}

	if _, err := c.App.UpdateChannelMemberSchemeRoles(c.Params.ChannelId, c.Params.UserId, schemeRoles.SchemeGuest, schemeRoles.SchemeUser, schemeRoles.SchemeAdmin); err != nil {
		c.Err = err
		return
	}

	auditRec.Success()

	ReturnStatusOK(w)
}

func updateChannelMemberNotifyProps(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireChannelId().RequireUserId()
	if c.Err != nil {
		return
	}

	props := model.MapFromJson(r.Body)
	if props == nil {
		c.SetInvalidParam("notify_props")
		return
	}

	auditRec := c.MakeAuditRecord("updateChannelMemberNotifyProps", audit.Fail)
	defer c.LogAuditRec(auditRec)
	auditRec.AddMeta("channel_id", c.Params.ChannelId)
	auditRec.AddMeta("props", props)

	if !c.App.SessionHasPermissionToUser(*c.AppContext.Session(), c.Params.UserId) {
		c.SetPermissionError(model.PermissionEditOtherUsers)
		return
	}

	_, err := c.App.UpdateChannelMemberNotifyProps(props, c.Params.ChannelId, c.Params.UserId)
	if err != nil {
		c.Err = err
		return
	}

	auditRec.Success()

	ReturnStatusOK(w)
}

func addChannelMember(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireChannelId()
	if c.Err != nil {
		return
	}

	props := model.StringInterfaceFromJson(r.Body)
	userId, ok := props["user_id"].(string)
	if !ok || !model.IsValidId(userId) {
		c.SetInvalidParam("user_id")
		return
	}

	member := &model.ChannelMember{
		ChannelId: c.Params.ChannelId,
		UserId:    userId,
	}

	postRootId, ok := props["post_root_id"].(string)
	if ok && postRootId != "" && !model.IsValidId(postRootId) {
		c.SetInvalidParam("post_root_id")
		return
	}

	if ok && len(postRootId) == 26 {
		rootPost, err := c.App.GetSinglePost(postRootId)
		if err != nil {
			c.Err = err
			return
		}
		if rootPost.ChannelId != member.ChannelId {
			c.SetInvalidParam("post_root_id")
			return
		}
	}

	channel, err := c.App.GetChannel(member.ChannelId)
	if err != nil {
		c.Err = err
		return
	}

	auditRec := c.MakeAuditRecord("addChannelMember", audit.Fail)
	defer c.LogAuditRec(auditRec)
	auditRec.AddMeta("channel", channel)

	if channel.Type == model.ChannelTypeDirect || channel.Type == model.ChannelTypeGroup {
		c.Err = model.NewAppError("addUserToChannel", "api.channel.add_user_to_channel.type.app_error", nil, "", http.StatusBadRequest)
		return
	}

	isNewMembership := false
	if _, err = c.App.GetChannelMember(context.Background(), member.ChannelId, member.UserId); err != nil {
		if err.Id == app.MissingChannelMemberError {
			isNewMembership = true
		} else {
			c.Err = err
			return
		}
	}

	isSelfAdd := member.UserId == c.AppContext.Session().UserId

	if channel.Type == model.ChannelTypeOpen {
		if isSelfAdd && isNewMembership {
			if !c.App.SessionHasPermissionToTeam(*c.AppContext.Session(), channel.TeamId, model.PermissionJoinPublicChannels) {
				c.SetPermissionError(model.PermissionJoinPublicChannels)
				return
			}
		} else if isSelfAdd && !isNewMembership {
			// nothing to do, since already in the channel
		} else if !isSelfAdd {
			if !c.App.SessionHasPermissionToChannel(*c.AppContext.Session(), channel.Id, model.PermissionManagePublicChannelMembers) {
				c.SetPermissionError(model.PermissionManagePublicChannelMembers)
				return
			}
		}
	}

	if channel.Type == model.ChannelTypePrivate {
		if isSelfAdd && isNewMembership {
			if !c.App.SessionHasPermissionToChannel(*c.AppContext.Session(), channel.Id, model.PermissionManagePrivateChannelMembers) {
				c.SetPermissionError(model.PermissionManagePrivateChannelMembers)
				return
			}
		} else if isSelfAdd && !isNewMembership {
			// nothing to do, since already in the channel
		} else if !isSelfAdd {
			if !c.App.SessionHasPermissionToChannel(*c.AppContext.Session(), channel.Id, model.PermissionManagePrivateChannelMembers) {
				c.SetPermissionError(model.PermissionManagePrivateChannelMembers)
				return
			}
		}
	}

	if channel.IsGroupConstrained() {
		nonMembers, err := c.App.FilterNonGroupChannelMembers([]string{member.UserId}, channel)
		if err != nil {
			if v, ok := err.(*model.AppError); ok {
				c.Err = v
			} else {
				c.Err = model.NewAppError("addChannelMember", "api.channel.add_members.error", nil, err.Error(), http.StatusBadRequest)
			}
			return
		}
		if len(nonMembers) > 0 {
			c.Err = model.NewAppError("addChannelMember", "api.channel.add_members.user_denied", map[string]interface{}{"UserIDs": nonMembers}, "", http.StatusBadRequest)
			return
		}
	}

	cm, err := c.App.AddChannelMember(c.AppContext, member.UserId, channel, app.ChannelMemberOpts{
		UserRequestorID: c.AppContext.Session().UserId,
		PostRootID:      postRootId,
	})
	if err != nil {
		c.Err = err
		return
	}

	auditRec.Success()
	auditRec.AddMeta("add_user_id", cm.UserId)
	c.LogAudit("name=" + channel.Name + " user_id=" + cm.UserId)

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(cm); err != nil {
		mlog.Warn("Error while writing response", mlog.Err(err))
	}
}

func removeChannelMember(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireChannelId().RequireUserId()
	if c.Err != nil {
		return
	}

	channel, err := c.App.GetChannel(c.Params.ChannelId)
	if err != nil {
		c.Err = err
		return
	}

	user, err := c.App.GetUser(c.Params.UserId)
	if err != nil {
		c.Err = err
		return
	}

	auditRec := c.MakeAuditRecord("removeChannelMember", audit.Fail)
	defer c.LogAuditRec(auditRec)
	auditRec.AddMeta("channel", channel)
	auditRec.AddMeta("remove_user_id", user.Id)

	if !(channel.Type == model.ChannelTypeOpen || channel.Type == model.ChannelTypePrivate) {
		c.Err = model.NewAppError("removeChannelMember", "api.channel.remove_channel_member.type.app_error", nil, "", http.StatusBadRequest)
		return
	}

	if channel.IsGroupConstrained() && (c.Params.UserId != c.AppContext.Session().UserId) && !user.IsBot {
		c.Err = model.NewAppError("removeChannelMember", "api.channel.remove_member.group_constrained.app_error", nil, "", http.StatusBadRequest)
		return
	}

	if c.Params.UserId != c.AppContext.Session().UserId {
		if channel.Type == model.ChannelTypeOpen && !c.App.SessionHasPermissionToChannel(*c.AppContext.Session(), channel.Id, model.PermissionManagePublicChannelMembers) {
			c.SetPermissionError(model.PermissionManagePublicChannelMembers)
			return
		}

		if channel.Type == model.ChannelTypePrivate && !c.App.SessionHasPermissionToChannel(*c.AppContext.Session(), channel.Id, model.PermissionManagePrivateChannelMembers) {
			c.SetPermissionError(model.PermissionManagePrivateChannelMembers)
			return
		}
	}

	if err = c.App.RemoveUserFromChannel(c.AppContext, c.Params.UserId, c.AppContext.Session().UserId, channel); err != nil {
		c.Err = err
		return
	}

	auditRec.Success()
	c.LogAudit("name=" + channel.Name + " user_id=" + c.Params.UserId)

	ReturnStatusOK(w)
}

func updateChannelScheme(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireChannelId()
	if c.Err != nil {
		return
	}

	schemeID := model.SchemeIDFromJson(r.Body)
	if schemeID == nil || !model.IsValidId(*schemeID) {
		c.SetInvalidParam("scheme_id")
		return
	}

	auditRec := c.MakeAuditRecord("updateChannelScheme", audit.Fail)
	defer c.LogAuditRec(auditRec)
	auditRec.AddMeta("new_scheme_id", schemeID)

	if c.App.Srv().License() == nil {
		c.Err = model.NewAppError("Api4.UpdateChannelScheme", "api.channel.update_channel_scheme.license.error", nil, "", http.StatusNotImplemented)
		return
	}

	if !c.App.SessionHasPermissionTo(*c.AppContext.Session(), model.PermissionManageSystem) {
		c.SetPermissionError(model.PermissionManageSystem)
		return
	}

	scheme, err := c.App.GetScheme(*schemeID)
	if err != nil {
		c.Err = err
		return
	}

	if scheme.Scope != model.SchemeScopeChannel {
		c.Err = model.NewAppError("Api4.UpdateChannelScheme", "api.channel.update_channel_scheme.scheme_scope.error", nil, "", http.StatusBadRequest)
		return
	}

	channel, err := c.App.GetChannel(c.Params.ChannelId)
	if err != nil {
		c.Err = err
		return
	}

	auditRec.AddMeta("channel", channel)
	auditRec.AddMeta("old_scheme_id", channel.SchemeId)

	channel.SchemeId = &scheme.Id

	_, err = c.App.UpdateChannelScheme(channel)
	if err != nil {
		c.Err = err
		return
	}

	auditRec.Success()

	ReturnStatusOK(w)
}

func channelMembersMinusGroupMembers(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireChannelId()
	if c.Err != nil {
		return
	}

	groupIDsParam := groupIDsQueryParamRegex.ReplaceAllString(c.Params.GroupIDs, "")

	if len(groupIDsParam) < 26 {
		c.SetInvalidParam("group_ids")
		return
	}

	groupIDs := []string{}
	for _, gid := range strings.Split(c.Params.GroupIDs, ",") {
		if !model.IsValidId(gid) {
			c.SetInvalidParam("group_ids")
			return
		}
		groupIDs = append(groupIDs, gid)
	}

	if !c.App.SessionHasPermissionTo(*c.AppContext.Session(), model.PermissionSysconsoleReadUserManagementChannels) {
		c.SetPermissionError(model.PermissionSysconsoleReadUserManagementChannels)
		return
	}

	users, totalCount, err := c.App.ChannelMembersMinusGroupMembers(
		c.Params.ChannelId,
		groupIDs,
		c.Params.Page,
		c.Params.PerPage,
	)
	if err != nil {
		c.Err = err
		return
	}

	b, marshalErr := json.Marshal(&model.UsersWithGroupsAndCount{
		Users: users,
		Count: totalCount,
	})
	if marshalErr != nil {
		c.Err = model.NewAppError("Api4.channelMembersMinusGroupMembers", "api.marshal_error", nil, marshalErr.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(b)
}

func channelMemberCountsByGroup(c *Context, w http.ResponseWriter, r *http.Request) {
	if c.App.Srv().License() == nil {
		c.Err = model.NewAppError("Api4.channelMemberCountsByGroup", "api.channel.channel_member_counts_by_group.license.error", nil, "", http.StatusNotImplemented)
		return
	}

	c.RequireChannelId()
	if c.Err != nil {
		return
	}

	if !c.App.SessionHasPermissionToChannel(*c.AppContext.Session(), c.Params.ChannelId, model.PermissionReadChannel) {
		c.SetPermissionError(model.PermissionReadChannel)
		return
	}

	includeTimezones := r.URL.Query().Get("include_timezones") == "true"

	channelMemberCounts, err := c.App.GetMemberCountsByGroup(app.WithMaster(context.Background()), c.Params.ChannelId, includeTimezones)
	if err != nil {
		c.Err = err
		return
	}

	b, marshalErr := json.Marshal(channelMemberCounts)
	if marshalErr != nil {
		c.Err = model.NewAppError("Api4.channelMemberCountsByGroup", "api.marshal_error", nil, marshalErr.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(b)
}

func getChannelModerations(c *Context, w http.ResponseWriter, r *http.Request) {
	if c.App.Srv().License() == nil {
		c.Err = model.NewAppError("Api4.GetChannelModerations", "api.channel.get_channel_moderations.license.error", nil, "", http.StatusNotImplemented)
		return
	}

	c.RequireChannelId()
	if c.Err != nil {
		return
	}

	if !c.App.SessionHasPermissionTo(*c.AppContext.Session(), model.PermissionSysconsoleReadUserManagementChannels) {
		c.SetPermissionError(model.PermissionSysconsoleReadUserManagementChannels)
		return
	}

	channel, err := c.App.GetChannel(c.Params.ChannelId)
	if err != nil {
		c.Err = err
		return
	}

	channelModerations, err := c.App.GetChannelModerationsForChannel(channel)
	if err != nil {
		c.Err = err
		return
	}

	b, marshalErr := json.Marshal(channelModerations)
	if marshalErr != nil {
		c.Err = model.NewAppError("Api4.getChannelModerations", "api.marshal_error", nil, marshalErr.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(b)
}

func patchChannelModerations(c *Context, w http.ResponseWriter, r *http.Request) {
	if c.App.Srv().License() == nil {
		c.Err = model.NewAppError("Api4.patchChannelModerations", "api.channel.patch_channel_moderations.license.error", nil, "", http.StatusNotImplemented)
		return
	}

	c.RequireChannelId()
	if c.Err != nil {
		return
	}

	auditRec := c.MakeAuditRecord("patchChannelModerations", audit.Fail)
	defer c.LogAuditRec(auditRec)

	if !c.App.SessionHasPermissionTo(*c.AppContext.Session(), model.PermissionSysconsoleWriteUserManagementChannels) {
		c.SetPermissionError(model.PermissionSysconsoleWriteUserManagementChannels)
		return
	}

	channel, appErr := c.App.GetChannel(c.Params.ChannelId)
	if appErr != nil {
		c.Err = appErr
		return
	}
	auditRec.AddMeta("channel", channel)

	var channelModerationsPatch []*model.ChannelModerationPatch
	err := json.NewDecoder(r.Body).Decode(&channelModerationsPatch)
	if err != nil {
		c.Err = model.NewAppError("Api4.patchChannelModerations", "api.marshal_error", nil, err.Error(), http.StatusInternalServerError)
		return
	}

	channelModerations, appErr := c.App.PatchChannelModerationsForChannel(channel, channelModerationsPatch)
	if appErr != nil {
		c.Err = appErr
		return
	}
	auditRec.AddMeta("patch", channelModerationsPatch)

	b, marshalErr := json.Marshal(channelModerations)
	if marshalErr != nil {
		c.Err = model.NewAppError("Api4.patchChannelModerations", "api.marshal_error", nil, marshalErr.Error(), http.StatusInternalServerError)
		return
	}

	auditRec.Success()
	w.Write(b)
}

func moveChannel(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireChannelId()
	if c.Err != nil {
		return
	}

	channel, err := c.App.GetChannel(c.Params.ChannelId)
	if err != nil {
		c.Err = err
		return
	}

	props := model.StringInterfaceFromJson(r.Body)
	teamId, ok := props["team_id"].(string)
	if !ok {
		c.SetInvalidParam("team_id")
		return
	}

	force, ok := props["force"].(bool)
	if !ok {
		c.SetInvalidParam("force")
		return
	}

	team, err := c.App.GetTeam(teamId)
	if err != nil {
		c.Err = err
		return
	}

	auditRec := c.MakeAuditRecord("moveChannel", audit.Fail)
	defer c.LogAuditRec(auditRec)
	auditRec.AddMeta("channel_id", channel.Id)
	auditRec.AddMeta("channel_name", channel.Name)
	auditRec.AddMeta("team_id", team.Id)
	auditRec.AddMeta("team_name", team.Name)

	if channel.Type == model.ChannelTypeDirect || channel.Type == model.ChannelTypeGroup {
		c.Err = model.NewAppError("moveChannel", "api.channel.move_channel.type.invalid", nil, "", http.StatusForbidden)
		return
	}

	if !c.App.SessionHasPermissionTo(*c.AppContext.Session(), model.PermissionManageSystem) {
		c.SetPermissionError(model.PermissionManageSystem)
		return
	}

	user, err := c.App.GetUser(c.AppContext.Session().UserId)
	if err != nil {
		c.Err = err
		return
	}

	err = c.App.RemoveAllDeactivatedMembersFromChannel(channel)
	if err != nil {
		c.Err = err
		return
	}

	if force {
		err = c.App.RemoveUsersFromChannelNotMemberOfTeam(c.AppContext, user, channel, team)
		if err != nil {
			c.Err = err
			return
		}
	}

	err = c.App.MoveChannel(c.AppContext, team, channel, user)
	if err != nil {
		c.Err = err
		return
	}

	auditRec.Success()
	c.LogAudit("channel=" + channel.Name)
	c.LogAudit("team=" + team.Name)

	if err := json.NewEncoder(w).Encode(channel); err != nil {
		mlog.Warn("Error while writing response", mlog.Err(err))
	}
}
