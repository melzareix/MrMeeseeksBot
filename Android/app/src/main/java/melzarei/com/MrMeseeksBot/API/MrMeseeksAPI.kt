package melzarei.com.MrMeseeksBot.API

import melzarei.com.MrMeseeksBot.Models.ChatRequest
import melzarei.com.MrMeseeksBot.Models.ChatResponse
import melzarei.com.MrMeseeksBot.Models.WelcomeResponse
import retrofit2.Call
import retrofit2.http.Body
import retrofit2.http.GET
import retrofit2.http.Header
import retrofit2.http.POST

interface MrMeseeksAPI {
    @GET("/welcome")
    /*
    * Get Welcome page to generate UUID.
    * */
    fun getUUID(): Call<WelcomeResponse>

    /*
    * Send Chat Message From User.
    * */
    @POST("/chat")
    fun sendChat(@Header("Authorization") content: String, @Body chatRequest: ChatRequest): Call<ChatResponse>
}