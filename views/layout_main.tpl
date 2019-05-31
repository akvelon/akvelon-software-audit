<!DOCTYPE html>

<html lang="en">
    <head>
        <title>Akvelon Software Audit | Scalable compliance and security audit for modern development</title>
        <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
        <link rel="stylesheet" href="/static/css/main.css">
    </head>
    <body>

    {{.Header}}

    <section class="main-content">
        <div class="landing-content landing-content"></div>
        <div class="container">
            <h2 class="subtitle">
                Enter the <strong>link</strong> of the GitHub repository to analyze:
            </h2>
            {{.LayoutContent}}
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
                    <tr>
                        <td class="table-link">Sorry, you have no reports generated yet.</td>
                    </tr>
                </tbody>
            </table>
        </div>
    </section>

    {{.Footer}}

    </body>
</html>