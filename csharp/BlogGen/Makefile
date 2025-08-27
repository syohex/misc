DOTNET_VERSION=$(strip $(shell perl -wne 'm/<TargetFramework>([^<]+)</ and print $$1' BlogGen.csproj))
LINUX_BINARY=bin/Release/$(DOTNET_VERSION)/linux-x64/publish/BlogGen
INSTALL_DIR=~/local/bloggen/bin

all: $(LINUX_BINARY)

$(LINUX_BINARY): *.cs
	dotnet publish -c Release -r linux-x64 --self-contained true -p:PublishSingleFile=true

install: $(LINUX_BINARY)
	rm -rf $(INSTALL_DIR)
	install -d $(INSTALL_DIR)
	cp $(LINUX_BINARY) $(INSTALL_DIR)/bloggen
	cp -r bin/Release/$(DOTNET_VERSION)/linux-x64/.playwright $(INSTALL_DIR)

.PHONY: clean
clean:
	dotnet clean
	rm -rf obj/ bin/
