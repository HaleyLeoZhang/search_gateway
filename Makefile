all: job

api:
	@clear
	@echo "App API debug is loading"
	@go run ./api/build/main.go -conf=./api/build/app.yaml

job:
	@clear
	@echo "App JOB debug is loading"
	@go run ./job/build/main.go -conf=./job/build/app.yaml

build:
	@clear
	@echo "App API is creating. Please wait ..."
	@rm -rf ./app_api
	@echo "App API compiling ..."
	@go build -o ./app_api -v ./api/build/main.go
	@echo "App API is created"

run:
	@#clear
	@echo "APP is loading. Please wait ..."
	@./app_api  -conf=./api/build/app.yaml

# 编译并运行项目当前 docker 内用
cr:
	@make build
	@make run

# 运行API项目
docker:
	@mkdir -p /tmp/comic_api
	@make -is docker_network
	@docker-compose down # 删除老的镜像
	@docker-compose up -d # 启动 docker-compile 编排 

docker_network:
	@docker network create --subnet=172.38.0.0/24  network_puppeteer_go # 创建 docker 网卡


ini:
	@clear
	@cp ./api/build/app.example.yaml ./api/build/app.yaml
	@echo "API config.yaml initial success"

tool:
	@clear
	@go vet ./...; true
	@gofmt -w .

clean:
	@clear
	@echo "Remove Old Apps ... "
	@rm -rf ./app_api
	@#rm -rf ./app_job # 2020-11-29 计划中
	@#rm -rf ./app_admin # 2020-11-14 暂无计划
	@go clean -i .
	@echo "Remove Old Apps --- Done"

test:
	@clear
	@echo "Test --- START"
	@# 全量测试---暂时不考虑
	@# go test -v ./...
	@# 指定测试
	@go test -v ./api/service/ -conf=../../api/build/app.yaml
	@echo "Test --- END"
