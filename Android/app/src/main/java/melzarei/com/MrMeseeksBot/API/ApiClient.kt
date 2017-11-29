package melzarei.com.MrMeseeksBot.API

import com.github.aurae.retrofit2.LoganSquareConverterFactory
import melzarei.com.MrMeseeksBot.Models.ChatError
import okhttp3.OkHttpClient
import okhttp3.ResponseBody
import okhttp3.logging.HttpLoggingInterceptor
import retrofit2.Retrofit
import java.util.concurrent.Executors

/*
* Singleton object used to send HTTP Requests to API.
* */

object ApiClient {
    private var retrofit: Retrofit
    val meseeksAPI: MrMeseeksAPI

    init {
        // Requests Logging
        val interceptor = HttpLoggingInterceptor()
        interceptor.level = HttpLoggingInterceptor.Level.BODY
        val client = OkHttpClient.Builder().addInterceptor(interceptor).build()
        retrofit = Retrofit
                .Builder()
                .baseUrl("http://35.193.174.93:3000")
                .client(client)
                .callbackExecutor(Executors.newSingleThreadExecutor())
                .addConverterFactory(LoganSquareConverterFactory.create())
                .build()
        meseeksAPI = retrofit.create(MrMeseeksAPI::class.java)
    }

    /*
    * Parse Retrofit Errors.
    * */
    fun parseErrors(response: ResponseBody): ChatError{
        val converter = ApiClient
                .retrofit
                .responseBodyConverter<ChatError>(ChatError::class.java, arrayOf<Annotation>())
        return converter.convert(response)
    }
}