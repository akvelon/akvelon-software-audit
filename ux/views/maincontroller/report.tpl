<section class="section container-results">
    <div class="container">
            <h1>Report for {{ .repoURL }}</h1>
    </div>

    <div class="container">
        <div class="columns">
            <table class="table">
                <colgroup>
                    <col width="65%">
                    <col width="22%">
                    <col width="13%">
                    <col width="10%">
                </colgroup>
                <thead>
                    <tr>
                        <th>File</th>
                        <th>Licence</th>
                        <th>Confidence</th>
                        <th>Size</th>
                    </tr>
                </thead>
                <tbody>
                    {{range $key, $val := .analyzeResult}}
                        <tr>
                            <td>{{$val.File}}</td>
                            <td>{{$val.License}}</td>
                            <td>{{$val.Confidence}}</td>
                            <td>{{$val.Size}}</td>
                        </tr>
                    {{end}}
                </tbody>
            </table>
        </div>
    </div>
</section>