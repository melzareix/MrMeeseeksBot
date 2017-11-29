package melzarei.com.MrMeseeksBot.Utils

import android.net.Uri
import android.support.customtabs.CustomTabsIntent
import android.support.v4.content.ContextCompat
import android.text.style.URLSpan
import android.view.View
import melzarei.com.MrMeseeksBot.R

class CustomTabsURLSpan(url: String) : URLSpan(url) {
    override fun onClick(widget: View) {
        val url = url
        val builder = CustomTabsIntent.Builder()
        builder.setToolbarColor(ContextCompat.getColor(widget.context, R.color.colorPrimary))
        val customTabsIntent = builder.build()
        try {
            customTabsIntent.launchUrl(widget.context, Uri.parse(url))
        } catch (e: Exception) {
            super.onClick(widget)
        }
    }
}