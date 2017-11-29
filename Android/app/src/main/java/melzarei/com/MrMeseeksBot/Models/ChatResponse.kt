package melzarei.com.MrMeseeksBot.Models

import com.bluelinelabs.logansquare.annotation.JsonField
import com.bluelinelabs.logansquare.annotation.JsonObject

/*
* ChatResponse JSON.
*/

@JsonObject(fieldDetectionPolicy = JsonObject.FieldDetectionPolicy.NONPRIVATE_FIELDS_AND_ACCESSORS)
class ChatResponse(
        @JsonField(name = arrayOf("anime_title"))
        var animeTitle: String = "",
        @JsonField(name = arrayOf("message"))
        var message: String = "",
        @JsonField(name = arrayOf("status"))
        var status: Boolean = false,
        @JsonField(name = arrayOf("code"))
        var code: Int = 0,
        @JsonField(name = arrayOf("imageURL"))
        var imageURL: String = "",
        @JsonField(name = arrayOf("google_calendar_authorized"))
        var googleCalendarAuthorized: Boolean = false
)

@JsonObject(fieldDetectionPolicy = JsonObject.FieldDetectionPolicy.NONPRIVATE_FIELDS_AND_ACCESSORS)
class WelcomeResponse(
        @JsonField(name = arrayOf("message"))
        var message: String = "",
        @JsonField(name = arrayOf("status"))
        var status: Boolean = false,
        @JsonField(name = arrayOf("code"))
        var code: Int = 0,
        @JsonField(name = arrayOf("uuid"))
        var uuid: String = ""
)