package melzarei.com.MrMeseeksBot.Utils

import android.graphics.Rect
import android.text.Spannable
import android.text.Spanned
import android.text.method.TransformationMethod
import android.text.style.URLSpan
import android.text.util.Linkify
import android.view.View
import android.widget.TextView


class LinkTransformationMethod : TransformationMethod {
    override fun getTransformation(source: CharSequence, view: View): CharSequence {
        if (view is TextView) {
            val textView = view as TextView
            Linkify.addLinks(textView, Linkify.WEB_URLS)
            val stringText = textView.text.toString()
            val text = textView.text as Spannable
            val spans = text.getSpans(0, textView.length(), URLSpan::class.java)
            for (i in spans.indices.reversed()) {
                val oldSpan = spans[i]
                text.removeSpan(oldSpan)
                val url = oldSpan.url
                val startIndex = stringText.indexOf(url)
                val lastIndex = startIndex + url.length
                text.setSpan(CustomTabsURLSpan(url), startIndex, lastIndex, Spanned.SPAN_EXCLUSIVE_EXCLUSIVE)
            }
            return text
        }
        return source
    }

    override fun onFocusChanged(view: View, sourceText: CharSequence, focused: Boolean,
                                direction: Int, previouslyFocusedRect: Rect) {

    }
}