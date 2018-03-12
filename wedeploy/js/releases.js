$.getJSON( "releases.json", function( data ) {
    var releases = [];
    $.each( data, function( id, release ) {
        releases.push(getReleaseHtml(id, release));
    });

    $( "<ul/>", {
        "class": "timeline timeline-center timeline-spacing-xl",
        html: releases.join( "" )
    }).appendTo( "#releases" );
});

function getReleaseHtml(index, release) {
    let reverseCssClass = "";

    if (release.channel === "stable") {
        reverseCssClass = "timeline-item-reverse"
    }

    let markup = `<li class="timeline-item ` + reverseCssClass + `">
    <div class="panel panel-secondary">
        <a aria-controls="panelCollapseTimelineSpacing` + index + `" aria-expanded="false" class="collapsed panel-header panel-header-link" data-toggle="collapse" href="#panelCollapseTimelineSpacing` + index + `" id="headingTimelineSpacing` + index + `" role="tab">
            <span class="panel-title">` + release.version + `</span>
            <div class="timeline-increment">
                <svg aria-hidden="true" class="lexicon-icon lexicon-icon-simulation-menu-closed">
                    <use xlink:href="/css/icons.svg#simulation-menu-closed" />
                </svg>
            </div>
            <div class="timeline-item-label">` + release.releaseDate + `</div>
        </a>
        <div aria-labelledby="headingTimelineSpacing` + index + `" class="collapse panel-collapse" id="panelCollapseTimelineSpacing` + index + `" role="tabpanel">
            <div id="panelBody` + index + `" class="panel-body">
                ` + getChangelog(release) + `
                ` + getDownloadLinks(release) + `
            </div>
        </div>
    </div>
</li>`;

    return markup;
}

function getChangelog(release) {
    let changelog = release.changelog;
    let enhacements = [];
    let breakings = [];
    let fixes = [];

    let changelogHtml = getChangelogMarkup(changelog.enhacements, enhacements, "Enhacements", "rocket");
    changelogHtml += getChangelogMarkup(changelog.breakings, breakings, "Breaking Changes", "skull");
    changelogHtml += getChangelogMarkup(changelog.fixes, fixes, "Fixes", "ant");

    return `<h3>
    <a aria-hidden="true">Changelog</a>. <a href="https://github.com/mdelapenya/lpn/releases/tag/` + release.version + `" target="_blank">See on Github</a>
</h3>` + changelogHtml;
}

function getChangelogMarkup(changelogElement, outputArray, label, emoji) {
    const header = `<h4>
    <a aria-hidden="true"> <i class="em em-` + emoji + `"></i> ` + label+ `</a>
</h4>`;

    if (!changelogElement || changelogElement.length == 0) {
        const nothing = `<div class="alert alert-info" role="alert">
    <span class="alert-indicator">
        <svg aria-hidden="true" class="lexicon-icon lexicon-icon-info-circle">
            <use xlink:href="/css/icons.svg#info-circle"></use>
        </svg>
    </span>
    <strong class="lead">Ups!</strong> Nothing interesting here.
</div>`;

        return header + nothing;
    }

    $.each(changelogElement, function(id, element) {
        outputArray.push(getDescription(element));
    });

    return  header + `<ul>` + outputArray.join( "" ) + `</ul>`;
}

function getDescription(change) {
    return `<li>` + change.description + `</li>`;
}

function getDownloadLinks(release) {
    if (release.equinox) {
        const equinoxUrl = 'https://dl.equinox.io/mdelapenya/lpn/stable';

        return `<a href='` + equinoxUrl + `' target='_blank'>Download from Equinox</a>`;
    }

    const header = `<h3><a aria-hidden="true"> <i class="em em-heart"></i> Downloads</a></h3>`;

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

    return header + `<ul>` + linksHtml + `</ul>`;
}