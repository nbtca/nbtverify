# SPDX-License-Identifier: GPL-2.0-only
#
# Copyright (C) 2024 nbtca

include $(TOPDIR)/rules.mk

PKG_NAME:=nbtverify
PKG_VERSION:=1.0.0
PKG_RELEASE:=1

PKG_SOURCE_PROTO:=git
PKG_SOURCE_URL:=https://github.com/nbtca/nbtverify.git
PKG_SOURCE_VERSION:=HEAD
PKG_MIRROR_HASH:=skip

PKG_LICENSE:=GPL-2.0
PKG_LICENSE_FILES:=LICENSE
PKG_MAINTAINER:=nbtca

PKG_BUILD_DEPENDS:=golang/host
PKG_BUILD_PARALLEL:=1
PKG_BUILD_FLAGS:=no-mips16

GO_PKG:=github.com/nbtca/nbtverify
GO_PKG_BUILD_PKG:=$(GO_PKG)
GO_PKG_LDFLAGS_X:=
GO_PKG_TAGS:=

include $(INCLUDE_DIR)/package.mk
include $(TOPDIR)/feeds/packages/lang/golang/golang-package.mk

define Package/nbtverify
  SECTION:=net
  CATEGORY:=Network
  TITLE:=NBT Campus Network Authentication Client
  URL:=https://github.com/nbtca/nbtverify
  DEPENDS:=$(GO_ARCH_DEPENDS)
endef

define Package/nbtverify/description
  NBT Campus Network Authentication Client for OpenWrt.
  Supports campus network portal authentication (卓智网络接入门户).
endef

define Package/nbtverify/conffiles
/etc/config/nbtverify
/etc/nbtverify/config.json
endef

define Package/nbtverify/install
	$(INSTALL_DIR) $(1)/usr/bin
	$(INSTALL_BIN) $(GO_PKG_BUILD_BIN_DIR)/nbtverify $(1)/usr/bin/
	$(INSTALL_DIR) $(1)/etc/nbtverify
	$(INSTALL_CONF) $(PKG_BUILD_DIR)/conf/config.json $(1)/etc/nbtverify/config.json
	$(INSTALL_DATA) $(PKG_BUILD_DIR)/conf/url.txt $(1)/etc/nbtverify/url.txt
endef

$(eval $(call GoBinPackage,nbtverify))
$(eval $(call BuildPackage,nbtverify))
