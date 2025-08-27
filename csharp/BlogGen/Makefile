DOTNET_VERSION=$(strip $(shell perl -wne 'm/<TargetFramework>([^<]+)</ and print $$1' BlogGen.csproj))
LINUX_BINARY=bin/Release/$(DOTNET_VERSION)/linux-x64/BlogGen

all: $(LINUX_BINARY)

$(LINUX_BINARY):
	dotnet publish -c Release -r linux-x64 --self-contained true /p:PublishSingleFile=true

install: $(LINUX_BINARY)
	cp $(LINUX_BINARY) ~/bin/bloggen

.PHONY: clean
clean:
	dotnet clean
