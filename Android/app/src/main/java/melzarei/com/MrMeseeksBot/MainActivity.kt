package melzarei.com.MrMeseeksBot

import android.media.MediaPlayer
import android.os.Bundle
import android.support.v7.app.AppCompatActivity
import android.util.Log
import android.widget.Toast
import com.squareup.picasso.Picasso
import com.stfalcon.chatkit.commons.ImageLoader
import com.stfalcon.chatkit.messages.MessageHolders
import com.stfalcon.chatkit.messages.MessagesListAdapter
import kotlinx.android.synthetic.main.activity_main.*
import melzarei.com.MrMeseeksBot.API.ApiClient
import melzarei.com.MrMeseeksBot.Models.*
import melzarei.com.MrMeseeksBot.Utils.IncomingMessageViewHolder
import retrofit2.Call
import retrofit2.Callback
import retrofit2.Response
import com.instabug.library.invocation.InstabugInvocationEvent
import com.instabug.library.Instabug




class MainActivity : AppCompatActivity() {

    private var currentUser: ChatUser? = null
    private lateinit var imageLoader: ImageLoader
    private lateinit var buttonClickPlayer: MediaPlayer
    private lateinit var newMessagePlayer: MediaPlayer
    lateinit var adapter: MessagesListAdapter<ChatListMessage>

    private fun showToast(text: String) {
        Toast
                .makeText(applicationContext, text, Toast.LENGTH_LONG)
                .show()
    }

    private fun sendChatMessage(dc: ChatRequest) {
        if (currentUser == null) {
            adapter.addToStart(ChatListMessage("Please wait, while I am acquiring your UserID."), true)
            return
        }
        appbar.subtitle = getString(R.string.alternate_enter_a_message)
        ApiClient
                .meseeksAPI.sendChat(currentUser!!.id, dc)
                .enqueue(object : Callback<ChatResponse> {
                    override fun onResponse(call: Call<ChatResponse>?, response: Response<ChatResponse>?) {
                        if (response!!.isSuccessful) {
                            val resp = response.body()!!
                            currentUser?.calendarAuthorized = resp.googleCalendarAuthorized
                            var scroll = true
                            runOnUiThread {
                                appbar.subtitle = getString(R.string.default_appbar_value)
                                newMessagePlayer.start()
                                if (resp.animeTitle != "") {
                                    adapter.addToStart(ChatListMessage(resp.animeTitle), scroll)
                                    scroll = false
                                }
                                if (resp.imageURL != "") {
                                    adapter.addToStart(ChatListMessage("", resp.imageURL), true)
                                    scroll = false
                                }
                            }
                            Thread.sleep(5) // Workaround for scrollAnimation
                            runOnUiThread {
                                adapter.addToStart(ChatListMessage(resp.message), scroll)
                            }

                        } else {
                            val error = ApiClient.parseErrors(response.errorBody()!!)
                            runOnUiThread {
                                appbar.subtitle = getString(R.string.default_appbar_value)
                                if (!buttonClickPlayer.isPlaying)
                                    newMessagePlayer.start()
                                adapter.addToStart(ChatListMessage(error.message), true)
                            }
                        }
                    }

                    override fun onFailure(call: Call<ChatResponse>?, t: Throwable?) {
                        val errorMessageText = t?.message ?: "Failed to connect to server."
                        runOnUiThread {
                            showToast(errorMessageText + " Retrying in 5 seconds.")
                        }
                        Thread.sleep(5000)
                        runOnUiThread {
                            sendChatMessage(dc)
                        }
                    }

                })
    }

    private fun handleWelcome() {
        ApiClient
                .meseeksAPI
                .getUUID()
                .enqueue(object : Callback<WelcomeResponse> {
                    override fun onFailure(call: Call<WelcomeResponse>?, t: Throwable?) {
                        val errorMessageText = t?.message ?: "Failed to connect to server."
                        runOnUiThread {
                            appbar.subtitle = getString(R.string.default_appbar_value)
                            showToast(errorMessageText + " Retrying in 5 seconds.")
                        }
                        Thread.sleep(5000)
                        runOnUiThread {
                            handleWelcome()
                        }
                    }

                    override fun onResponse(call: Call<WelcomeResponse>?, response: Response<WelcomeResponse>?) {
                        if (response!!.isSuccessful) {
                            currentUser = ChatUser(response.body()!!.uuid, response.body()!!.uuid)
                            runOnUiThread {
                                appbar.subtitle = getString(R.string.default_appbar_value)
                                adapter.addToStart(ChatListMessage(response.body()!!.message), true)
                            }
                        }
                    }

                })
    }

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)
        appbar.title = "Mr Meseeks"
        appbar.subtitle = getString(R.string.alternate_enter_a_message)

        // MediaPlayers
        buttonClickPlayer = MediaPlayer.create(this, R.raw.button_click)
        newMessagePlayer = MediaPlayer.create(this, R.raw.button_click)

        // ImageLoader
        imageLoader = ImageLoader { imageView, url -> Picasso.with(applicationContext).load(url).into(imageView) }

        // Setup Message Adapter
        val holdersConfig = MessageHolders()
        holdersConfig.setIncomingImageLayout(R.layout.custom_incoming_image_message)
        holdersConfig.setIncomingTextConfig(IncomingMessageViewHolder::class.java, R.layout.custom_incoming_message)
        adapter = MessagesListAdapter("User", holdersConfig, imageLoader)
        adapter.disableSelectionMode()
        messagesList.setAdapter(adapter)

        // GET UUID
        handleWelcome()

        // Attachment Button
        input.setAttachmentsListener {
            val authorized = currentUser?.calendarAuthorized ?: false
            if (authorized) {
                adapter.addToStart(ChatListMessage("Google Calendar Already Authorized."), true)
            } else {
                adapter.addToStart(ChatListMessage("Authorize Calendar", null, "User"), true)
                sendChatMessage(ChatRequest("Authorize Calendar"))
            }
        }

        // Input Button
        input.setInputListener({ data ->
            buttonClickPlayer.start()
            adapter.addToStart(ChatListMessage(data.toString(), null, "User"), true)
            val dc = ChatRequest(data.toString())
            sendChatMessage(dc)
            true
        })

        // InstaBug integration
        Instabug.Builder(application, "4d9210747d8131ca5d29bb2fe129a01e")
                .setInvocationEvent(InstabugInvocationEvent.SHAKE)
                .build()

    }
}
