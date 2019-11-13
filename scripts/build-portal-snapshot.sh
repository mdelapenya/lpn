#!/bin/bash

readonly RELEASE_FILE_URL="https://releases.liferay.com/portal/snapshot-master/latest/liferay-portal-tomcat-master.7z"

function build_docker_image {
	local docker_image_name
	local label_name

	if [[ ${RELEASE_FILE_NAME} == *-portal-* ]]
	then
		docker_image_name="portal"
		label_name="Liferay Portal"
	else
		echo "${RELEASE_FILE_NAME} is an unsupported release file name."

		exit 1
	fi

	if [[ ${RELEASE_FILE_URL%} == */snapshot-* ]]
	then
		docker_image_name=${docker_image_name}-snapshot
	fi

	if [[ ${RELEASE_FILE_URL} == http://release* ]]
	then
		docker_image_name=${docker_image_name}-snapshot
	fi

	local release_version=${RELEASE_FILE_URL%/*}

	release_version=${release_version##*/}

	if [[ ${RELEASE_FILE_URL} == http://release* ]]
	then
		release_version=${RELEASE_FILE_URL#*tomcat-}
		release_version=${release_version%.*}
	fi

	local label_version=${release_version}

	if [[ ${RELEASE_FILE_URL%} == */snapshot-* ]]
	then
		local release_branch=${RELEASE_FILE_URL%/*}

		release_branch=${release_branch%/*}
		release_branch=${release_branch%-private*}
		release_branch=${release_branch##*-}

		local release_hash=$(cat ${TEMP_DIR}/liferay/.githash)

		release_hash=${release_hash:0:7}

		if [[ ${release_branch} == master ]]
		then
			label_version="Master Snapshot on ${label_version} at ${release_hash}"
		else
			label_version="${release_branch} Snapshot on ${label_version} at ${release_hash}"
		fi
	fi

	DOCKER_IMAGE_TAGS=()
    DOCKER_IMAGE_TAGS+=("mdelapenya/${docker_image_name}:${release_branch}-${release_version}-${release_hash}")
    DOCKER_IMAGE_TAGS+=("mdelapenya/${docker_image_name}:${release_branch}-$(date "${CURRENT_DATE}" "+%Y%m%d")")
    DOCKER_IMAGE_TAGS+=("mdelapenya/${docker_image_name}:${release_branch}")

	docker build \
		--build-arg LABEL_BUILD_DATE=$(date "${CURRENT_DATE}" "+%Y-%m-%dT%H:%M:%SZ") \
		--build-arg LABEL_NAME="${label_name}" \
		--build-arg LABEL_VCS_REF=$(git rev-parse HEAD) \
		--build-arg LABEL_VCS_URL="https://github.com/liferay/liferay-docker" \
		--build-arg LABEL_VERSION="${label_version}" \
		$(get_docker_image_tags_args ${DOCKER_IMAGE_TAGS[@]}) \
		${TEMP_DIR}
}

function clone_build_scripts {
    git clone https://github.com/liferay/liferay-docker.git
}

function main {
    clone_build_scripts || true

    cd liferay-docker
    source ./_common.sh

	make_temp_directory

	prepare_temp_directory ${RELEASE_FILE_URL}

	prepare_tomcat

	build_docker_image

	push_docker_images ${1}

	clean_up_temp_directory
}

function prepare_temp_directory {
	RELEASE_FILE_NAME=${1##*/}

	release_file_url=${1}

	local release_dir=${1%/*}

	release_dir=${release_dir#*com/}
	release_dir=${release_dir#*com/}
	release_dir=${release_dir#*liferay-release-tool/}
	release_dir=${release_dir#*private/ee/}
	release_dir=releases/${release_dir}

	if [ ! -e ${release_dir}/${RELEASE_FILE_NAME} ]
	then
		echo ""
		echo "Downloading ${release_file_url} to ${release_dir}."
		echo ""

		mkdir -p ${release_dir}

		curl -f -o ${release_dir}/${RELEASE_FILE_NAME} ${release_file_url} || exit 2
	fi

	if [[ ${RELEASE_FILE_NAME} == *.7z ]]
	then
		7z x -O${TEMP_DIR} ${release_dir}/${RELEASE_FILE_NAME} || exit 3
	else
		unzip -q ${release_dir}/${RELEASE_FILE_NAME} -d ${TEMP_DIR}  || exit 3
	fi

	mv ${TEMP_DIR}/liferay-* ${TEMP_DIR}/liferay
}

main "$@"