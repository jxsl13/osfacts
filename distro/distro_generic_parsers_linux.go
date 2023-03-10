package distro

// Uses the provided distribution name or its alias as the distribution name.
// Tries to search through te whole text file for one or more semantic version string matches.
// The longest and highest semantic version is used.
func parserFindSemanticVersion(dist distribution, filePath, fileContent string, osInfo *Info) error {
	distVersion, err := findSemanticVersionString(fileContent)
	if err != nil {
		return err
	}

	osInfo.update(dist.InfoName(), distVersion)
	return nil
}

// Expects a key value file format (.env)
// Uses the distribution name define din the distribution objects and solely detects the os version
// base donthe provided keys. In cas eno keys wer eprovided, then all keys are searched through for
// semantic versions.
func parserFindEnvSemanticVersionKeys(keys ...string) fileParseFunc {
	return func(dist distribution, filePath, fileContent string, osInfo *Info) error {
		distVersion, err := findEnvSemanticVersion(fileContent, keys...)
		if err != nil {
			return err
		}
		osInfo.update(dist.InfoName(), distVersion)
		return nil
	}
}

// Expects a key value file format (.env)
// looks for the name key and used it and looks through on eor all keys for a semantic versio string.
func parserFindEnvNameAndSemanticVersionKeys(distNameKey string, versionKeys ...string) fileParseFunc {
	return func(dist distribution, filePath, fileContent string, osInfo *Info) error {
		envMap, err := getEnvMap(fileContent)
		if err != nil {
			return err
		}

		distName, err := getKey(envMap, distNameKey)
		if err != nil {
			return err
		}
		distVersion, err := findEnvSemanticVersion(fileContent, versionKeys...)
		if err != nil {
			return err
		}

		osInfo.update(distName, distVersion)
		return nil
	}
}

func parserFindEnvVersionKey(versionKeys string) fileParseFunc {
	return func(dist distribution, filePath, fileContent string, osInfo *Info) error {
		envMap, err := getEnvMap(fileContent)
		if err != nil {
			return err
		}

		distVersion, err := getKey(envMap, versionKeys)
		if err != nil {
			return err
		}

		osInfo.update(dist.InfoName(), distVersion)
		return nil
	}
}
