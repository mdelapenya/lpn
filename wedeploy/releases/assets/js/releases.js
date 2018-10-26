function getStableRelease() {
    $.getJSON( "releases.json", function( data ) {
        var releases = [];
        $.each( data, function( id, release ) {
            if (release.latestStable && release.latestStable === true) {
                releases.push(getReleaseHtml(id, release));

                return;
            }
        });
    
        $( "<div/>", {
            "class": "accordion",
            "id": "accordion",
            html: releases.join( "" )
        }).appendTo( "#releases" );
    });
}

function getReleases() {
    $.getJSON( "releases.json", function( data ) {
        var releases = [];
        $.each( data, function( id, release ) {
            releases.push(getReleaseHtml(id, release));
        });
    
        $( "<div/>", {
            "class": "accordion",
            "id": "accordion",
            html: releases.join( "" )
        }).appendTo( "#releases" );
    });
}

function getReleaseHtml(index, release) {
    let timeLineItemLabel = "In development";
    let collapseShow = (index == 0) ? "show" : ""
    let latestStable = ""
    let textWarning = "text-warning"

    if (release.channel === "stable") {
        timeLineItemLabel = release.releaseDate;
        textWarning = ""
    }

    if (release.latestStable && release.latestStable === true) {
        latestStable = ` &nbsp;<span class="badge badge-success">Latest Stable</span>`
    }

    let markup = `<div class="card">
    <a class="card-header" id="heading` + index + `" data-toggle="collapse" data-target="#collapse` + index + `" aria-expanded="true" aria-controls="collapse` + index + `">
        <div class="row">
            <div class="col-md-2 col-sm-12 version-date small d-flex align-self-center ` + textWarning + `">
                <span>`+ timeLineItemLabel + `</span>
            </div>
            <div class="col-md-10 col-sm-12">
                <span class="version-title">Version ` + release.version + `</span>` + latestStable + `
                <svg class="lexicon-icon lexicon-icon-angle-down" viewBox="0 0 512 512">
                    <path class="lexicon-icon-outline" d="M256 384c6.936-0.22 13.81-2.973 19.111-8.272l227.221-227.221c11.058-11.026 11.058-28.941 0-39.999-11.026-11.058-28.94-11.058-39.999 0l-206.333 206.333c0 0-206.333-206.333-206.333-206.333-11.059-11.058-28.973-11.058-39.999 0-11.059 11.058-11.059 28.972 0 39.999l227.221 227.221c5.3 5.3 12.174 8.053 19.111 8.272z"></path>
                </svg>
            </div>
        </div>
    </a>
    <div id="collapse` + index + `" class="collapse ` + collapseShow + `" aria-labelledby="heading` + index + `" data-parent="#accordion">
        <div class="card-body">
            <div class="row">
                <div class="col-md-5 offset-md-2 col-sm-12">
                ` + getChangelog(release) + `
                </div>
                <div class="col-md-4  offset-md-1 col-sm-12">` +
                    getDownloadLinks(release) +
                `</div>
            </div>
        </div>
    </div>
</div>`;

    return markup;
}

function getChangelog(release) {
    let changelog = release.changelog;
    let enhacements = [];
    let breakings = [];
    let fixes = [];

    let changelogHtml = getChangelogMarkup(changelog.enhacements, enhacements, "Enhacements");
    changelogHtml += getChangelogMarkup(changelog.breakings, breakings, "Breaking Changes");
    changelogHtml += getChangelogMarkup(changelog.fixes, fixes, "Fixes");
    changelogHtml += `<a href="https://github.com/mdelapenya/lpn" class="">See it on Github</a>`;

    return changelogHtml;
}

function getChangelogMarkup(changelogElement, outputArray, label) {    
    if (!changelogElement || changelogElement.length == 0) {
        return '';
    }

    const header = `<h5 class="list-section">` + label+ `</h5>`;

    $.each(changelogElement, function(id, element) {
        outputArray.push(getDescription(element));
    });

    return  header + `<ul class="release-list">` + outputArray.join( "" ) + `</ul>`;
}

function getDescription(change) {
    return `<li>` + change.description + `</li>`;
}

function getDownloadLinks(release) {
    if (release.equinox) {
        const equinoxUrl = 'https://dl.equinox.io/mdelapenya/lpn/stable';

        return `<a href='` + equinoxUrl + `' target='_blank'>Download from Equinox</a>`;
    }

    const header = `<h5 class="list-section">Downloads</h5>`;

    const oss = ['darwin', 'linux', 'windows'];
    const platforms = ['386', 'amd64'];

    let basePath = `/bin/` + release.channel + `/` + release.version;

    let linksHtml = '';

    $.each(oss, function(id, os) {
        let extension = '';

        if (os === 'windows') {
            extension = '.exe';
        }

        $.each(platforms, function(id, platform) {
            let url = basePath + `/` + os + `/` + platform + `/lpn` + extension;

            linksHtml += `<li>
    <a href='` + url + `' target='_blank'>` + os + ` - ` + platform + `</a>
</li>`
        });
    });

    return header + `<ul class="release-list links">` + linksHtml + `</ul>`;
}