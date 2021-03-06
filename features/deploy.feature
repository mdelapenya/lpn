Feature: Deploy command
  As a newcomer to lpn
  I want to be able to deploy files and directories to the container created by the tool

  Scenario Outline: Deploy one single file when container exists
    Given an empty file named "modules/a.jar"
    When I run `lpn run <type> -t <tag>`
    And I run `docker exec lpn-<type> mkdir -p <home>`
    And I run `lpn deploy <type> -f modules/a.jar`
    Then the output should contain:
    """
    File deployed successfully to deploy dir
    """
    And the output should contain:
    """
    file=modules/a.jar
    """
    And the output should contain:
    """
    deployDir=<home>
    """
    And I run `docker exec lpn-<type> ls -l <home> | grep "a.jar" | wc -l | xargs`
    And the output should contain:
    """
    1
    """
    And I run `lpn rm <type>`

  Examples:
    | type | tag | home |
    | ce      | 7.0.6-ga7 | /opt/liferay/deploy |
    | commerce | 1.1.1 | /opt/liferay/deploy |
    | dxp     | 7.0.10.8 | /opt/liferay/deploy |
    | nightly | master | /opt/liferay/deploy |
    | release | latest | /liferay/deploy |
    | release | 7-ce-ga5-tomcat-hsql | /usr/local/liferay-ce-portal-7.0-ga5/deploy |

  Scenario Outline: Deploy multiple files when container exists
    Given an empty file named "modules/a.jar"
    And an empty file named "modules/b.jar"
    When I run `lpn run <type> -t <tag>`
    And I run `docker exec lpn-<type> mkdir -p <home>`
    And I run `lpn deploy <type> -f modules/a.jar,modules/b.jar`
    Then the output should contain:
    """
    File deployed successfully to deploy dir
    """
    And the output should contain:
    """
    file=modules/a.jar
    """
    And the output should contain:
    """
    file=modules/b.jar
    """
    And the output should contain:
    """
    deployDir=<home>
    """
    And I run `docker exec lpn-<type> ls -l <home> | grep "a.jar" | wc -l | xargs`
    And the output should contain:
    """
    1
    """
    And I run `docker exec lpn-<type> ls -l <home> | grep "b.jar" | wc -l | xargs`
    And the output should contain:
    """
    1
    """
    And I run `lpn rm <type>`

  Examples:
    | type | tag | home |
    | ce      | 7.0.6-ga7 | /opt/liferay/deploy |
    | commerce | 1.1.1 | /opt/liferay/deploy |
    | dxp     | 7.0.10.8 | /opt/liferay/deploy |
    | nightly | master | /opt/liferay/deploy |
    | release | latest | /liferay/deploy |
    | release | 7-ce-ga5-tomcat-hsql | /usr/local/liferay-ce-portal-7.0-ga5/deploy |

  Scenario Outline: Deploy command with no flags
    Given I run `lpn run <type> -t <tag>`
    When I run `lpn deploy <type>`
    Then the output should contain:
    """
    Please pass a valid path to a file or to a directory as argument
    """
    And the exit status should be 1
    And I run `lpn rm <type>`

  Examples:
    | type    | tag |
    | ce      | 7.0.6-ga7 |
    | commerce | 1.1.1 |
    | dxp     | 7.0.10.8 |
    | nightly | master |
    | release | latest |

  Scenario Outline: Deploy a directory when container exists
    Given an empty file named "modules/a.jar"
    And an empty file named "modules/b.jar"
    And an empty file named "modules/c.jar"
    When I run `lpn run <type> -t <tag>`
    And I run `docker exec lpn-<type> mkdir -p <home>`
    And I run `lpn deploy <type> -d modules`
    Then the output should contain:
    """
    File deployed successfully to deploy dir
    """
    And the output should contain:
    """
    file=modules/a.jar
    """
    And I run `docker exec lpn-<type> ls -l <home> | grep "a.jar" | wc -l | xargs`
    And the output should contain:
    """
    1
    """
    And the output should contain:
    """
    file=modules/b.jar
    """
    And I run `docker exec lpn-<type> ls -l <home> | grep "b.jar" | wc -l | xargs`
    And the output should contain:
    """
    1
    """
    And the output should contain:
    """
    file=modules/c.jar
    """
    And I run `docker exec lpn-<type> ls -l <home> | grep "c.jar" | wc -l | xargs`
    And the output should contain:
    """
    1
    """
    And I run `lpn rm <type>`

    Examples:
    | type | tag | home |
    | ce      | 7.0.6-ga7 | /opt/liferay/deploy |
    | commerce | 1.1.1 | /opt/liferay/deploy |
    | dxp     | 7.0.10.8 | /opt/liferay/deploy |
    | nightly | master | /opt/liferay/deploy |
    | release | latest | /liferay/deploy |
    | release | 7-ce-ga5-tomcat-hsql | /usr/local/liferay-ce-portal-7.0-ga5/deploy |

  Scenario Outline: Deploy a directory skipping subdirectories when container exists
    Given an empty directory named "modules/skip1"
    And an empty directory named "modules/skip2"
    When I run `lpn run <type> -t <tag>`
    And I run `docker exec lpn-<type> mkdir -p <home>`
    And I run `lpn deploy <type> -d modules`
    Then the output should not contain:
    """
    File deployed successfully to deploy dir
    """
    And the output should noy contain:
    """
    file=modules/skip1
    """
    And I run `docker exec lpn-<type> ls -l <home> | grep "skip1" | wc -l | xargs`
    And the output should contain:
    """
    0
    """
    And the output should not contain:
    """
    File deployed successfully to deploy dir
    """
    And the output should not contain:
    """
    file=modules/skip2
    """
    And I run `docker exec lpn-<type> ls -l <home> | grep "skip2" | wc -l | xargs`
    And the output should contain:
    """
    0
    """
    And I run `lpn rm <type>`

    Examples:
    | type | tag | home |
    | ce      | 7.0.6-ga7 | /opt/liferay/deploy |
    | commerce | 1.1.1 | /opt/liferay/deploy |
    | dxp     | 7.0.10.8 | /opt/liferay/deploy |
    | nightly | master | /opt/liferay/deploy |
    | release | latest | /liferay/deploy |
    | release | 7-ce-ga5-tomcat-hsql | /usr/local/liferay-ce-portal-7.0-ga5/deploy |

  Scenario Outline: Deploy a directory skipping subdirectories when container does not exist
    Given an empty directory named "modules/skip1"
    And an empty directory named "modules/skip2"
    When I run `lpn rm <type>`
    And I run `lpn deploy <type> -d modules`
    Then the output should contain:
    """
    We could not find the container among the running containers
    """
    And the output should contain:
    """
    container=lpn-<type>
    """
    And the exit status should be 1

    Examples:
    | type     |
    | ce       |
    | commerce |
    | dxp      |
    | nightly  |
    | release  |
