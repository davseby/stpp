<script setup lang="ts">
import { RouterView } from "vue-router"
import { NGradientText, NNotificationProvider } from "naive-ui"
import TheNavbar from "@/components/TheNavbar.vue"
</script>

<template>
	<n-notification-provider placement="top-right">
		<div class="title-wrapper">
			<div class="title">
				<n-gradient-text type="success">
					Foodie
				</n-gradient-text>
			</div>
		</div>
		<div class="body-wrapper">
			<div class="divider"></div>
			<TheNavbar />
			<div class="content-wrapper">
				<RouterView />
			</div>
		</div>
	</n-notification-provider>
</template>

<style lang="scss" scoped>
@use "./assets/main.scss";

.title-wrapper {
	position: fixed;
	top: 0;
	z-index: 100;
	width: 100%;

	.title {
		font-size: 45px;
		text-align: center;
		font-weight: 600;
		background-color: white;
	}
}

.body-wrapper {
	margin: 72px 8px 0 8px;
	
	.content-wrapper {
		margin: 10px;
	}
}
</style>

