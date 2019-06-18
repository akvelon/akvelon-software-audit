<section class="main-content">
    <div class="landing-content landing-content"></div>
    <div class="container">
        <h2 class="subtitle">
            Enter the <strong>link</strong> of the GitHub repository to analyze:
        </h2>

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
        <form method="POST" action="/analyze" id="check_form">
                <div>
                    <p>
                    <input name="repo" type="text" class="input-box" placeholder="GitHub repo link goes here..."/>
                    </p>
                </div>
                <div>
                    <button type="submit" class="button btn" href="#" role="button">Generate Report</button>
                </div>
        </form>
    </div>
</section>

<section class="section">
    <div class="container">
        <table class="table">
            <thead>
                <tr>
                <th>Recently Generated</th>
                </tr>
            </thead>
            <tbody>
                {{ range $key, $val := .Recent }}
                        <tr>
                            <td class="table-link"><a href="/report/{{ $val }}">{{ $val }}</td>
                        </tr>
                {{ end }}
            </tbody>
        </table>
    </div>
</section>