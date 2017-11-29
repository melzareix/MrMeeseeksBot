package melzarei.com.MrMeseeksBot.Models

import com.bluelinelabs.logansquare.annotation.JsonField
import com.bluelinelabs.logansquare.annotation.JsonObject


/*
* Error Message JSON Model.
* */

@JsonObject(fieldDetectionPolicy = JsonObject.FieldDetectionPolicy.NONPRIVATE_FIELDS_AND_ACCESSORS)
class ChatError(
        @JsonField(name = arrayOf("message"))
        var message: String = "",
        @JsonField(name = arrayOf("status"))
        var status: Boolean = false,
        @JsonField(name = arrayOf("code"))
        var code: Int = 400

) {
        override fun toString(): String {
                return "ChatError(message='$message', status=$status, code=$code)"
        }
}