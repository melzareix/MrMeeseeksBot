package melzarei.com.MrMeseeksBot.Utils

import android.text.method.LinkMovementMethod
import android.view.View
import android.widget.TextView
import com.stfalcon.chatkit.messages.MessageHolders
import melzarei.com.MrMeseeksBot.Models.ChatListMessage

/**
 * Created by melzarei on 11/29/17.
 * mohamedelzarei@gmail.com
 */

class IncomingMessageViewHolder(itemView: View) :
        MessageHolders.IncomingTextMessageViewHolder<ChatListMessage>(itemView) {
    override fun configureLinksBehavior(text: TextView?) {
        super.configureLinksBehavior(text)
        text?.transformationMethod = LinkTransformationMethod()
        text?.movementMethod = LinkMovementMethod.getInstance()
    }
}