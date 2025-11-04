import MarkdownIt from 'markdown-it'
import mk from "@ruanyf/markdown-it-katex";
var hljs = require('highlight.js');
hljs.configure({
    lineNumbers: true
});
import 'highlight.js/styles/atom-one-dark.css';

export const md = MarkdownIt({
    // 在源码中启用 HTML 标签
    html: true,
    // 如果结果以 <pre ... 开头，内部包装器则会跳过。
    highlight: function(str, lang) {

        // 经过highlight.js处理后的html
        let preCode = ""
        try {
            if(lang && hljs.getLanguage(lang)){
                preCode = hljs.highlight(str, { language: lang }).value
            }else{
                preCode = hljs.highlightAuto(str).value;
            }
        } catch (err) {
            preCode = md.utils.escapeHtml(str);
        }

        const lines = preCode.split(/\n/).slice(0, -1)
        let _lines = lines.filter((it,i) => it !== '')

        // 添加自定义行号
        let html = _lines.map((item, index) => {
            return '<li class="line-li"><span class="line-numbers-rows"></span>' + item +
                '</li>'
        }).join('')
        html = '<ol style="padding: 0px 30px;">' + html + '</ol>'

        // 代码复制功能
        let htmlCode =
            `<div style="color: #888;border-radius: 0 0 5px 5px;">`

        htmlCode += `<div class="code-header">`
        htmlCode +=
            `${lang}<a class="copy-btn mk-copy-btn" style="cursor: pointer;">复制 </a>`
        htmlCode += `</div>`

        htmlCode +=
            `<pre class="hljs" style="padding:0 10px!important;margin-bottom:5px;overflow: auto;display: block;border-radius: 5px;"><code>${html}</code></pre>`;
        htmlCode += '</div>'
        return htmlCode
    }
})

md.use(mk, { "throwOnError": false, "errorColor": "#000000", "output": "mathml"})
