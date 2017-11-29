package melzarei.com.MrMeseeksBot.Models

import com.stfalcon.chatkit.commons.models.IUser

/*
* Chat User Class.
* */

class ChatUser(private val userName: String, private val userUUID: String,
               var calendarAuthorized: Boolean = false) : IUser {
    override fun getAvatar(): String {
        if (userName == "Meseeks")
            return "https://i.imgur.com/TQGbW2B.jpg" // Mr Meseeks avatar
        return ""
    }

    override fun getName(): String {
        return userName
    }

    override fun getId(): String {
        return userUUID
    }

    override fun toString(): String {
        return "ChatUser(userName='$userName', userUUID='$userUUID', calendarAuthorized=$calendarAuthorized)"
    }


}