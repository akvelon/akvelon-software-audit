<div id="notifications">
    {{if .flash.error }}
        <div class="notification is-error">
            {{.flash.error}}
        </div>
    {{end}}
    {{if .flash.success }}
        <div class="notification is-success">
            {{.flash.success}}
        </div>
    {{end}}
</div>