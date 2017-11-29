package melzarei.com.MrMeseeksBot.Models

import com.stfalcon.chatkit.commons.models.IMessage
import com.stfalcon.chatkit.commons.models.IUser
import com.stfalcon.chatkit.commons.models.MessageContentType
import java.util.*


/**
 * Created by melzarei on 11/26/17.
 * mohamedelzarei@gmail.com
 */

class ChatListMessage(val message: String, val imageURL: String? = null, val userName: String = "Meseeks", val messageId: String = "1") :
        IMessage, MessageContentType.Image {

    override fun getImageUrl(): String? {
        if (userName == "Meseeks")
            return imageURL
        return null
    }

    override fun getId(): String {
        return messageId
    }

    override fun getCreatedAt(): Date {
        return Calendar.getInstance().time
    }

    override fun getUser(): IUser {
        return ChatUser(userName, userName)
    }

    override fun getText(): String {
        return message
    }

}