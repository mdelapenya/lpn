$.getJSON( "releases.json", function( data ) {
    var releases = [];
    $.each( data, function( id, release ) {
        releases.push(getReleaseHtml(id, release));
    });

    $( "<ul/>", {
        "class": "timeline timeline-center timeline-spacing-xl",
        html: releases.join( "" )
    }).appendTo( "body" );
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
            <div class="timeline-item-label">` + release.commitDate + `</div>
        </a>
        <div aria-labelledby="headingTimelineSpacing` + index + `" class="collapse panel-collapse" id="panelCollapseTimelineSpacing` + index + `" role="tabpanel">
            <div id="panelBody` + index + `" class="panel-body">
                ` + getChangelog(index, release) + `
            </div>
        </div>
    </div>
</li>`;

    return markup;
}

function getChangelog(index, release) {
    let changelog = release.changelog;
    let enhacements = [];
    let breakings = [];
    let fixes = [];

    let changelogHtml = getChangelogMarkup(changelog.enhacements, enhacements, "Enhacements", "rocket");
    changelogHtml += getChangelogMarkup(changelog.breakings, breakings, "Breaking Changes", "skull");
    changelogHtml += getChangelogMarkup(changelog.fixes, fixes, "Fixes", "ant");

    return changelogHtml;
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