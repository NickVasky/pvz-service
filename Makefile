codegen_api_spec := ./_specs/swagger.yaml
codegen_dir := ./codegen
codegen_options := #--global-property models,apis
codegen_config := ./_specs/openapi_generator_config.yaml #./_specs/oapi_codegen_config.yaml


dto_generate:
# go tool oapi-codegen -config ${codegen_config} ${codegen_api_spec}
	openapi-generator-cli generate -g go-server -c ${codegen_config} -i ${codegen_api_spec} -o ${codegen_dir} ${codegen_options}
#cd ${codegen_dir}; go mod tidy

dto_clean:
	rm -rf ${codegen_dir}

dto_regenerate: dto_clean dto_generate
	echo "regen"