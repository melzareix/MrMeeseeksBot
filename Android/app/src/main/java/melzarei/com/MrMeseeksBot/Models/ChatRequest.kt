package melzarei.com.MrMeseeksBot.Models

import com.bluelinelabs.logansquare.annotation.JsonField
import com.bluelinelabs.logansquare.annotation.JsonObject

/*
* ChatRequest JSON.
*/

@JsonObject(fieldDetectionPolicy = JsonObject.FieldDetectionPolicy.NONPRIVATE_FIELDS_AND_ACCESSORS)
class ChatRequest(
        @JsonField(name = arrayOf("message"))
        var message: String = ""
)