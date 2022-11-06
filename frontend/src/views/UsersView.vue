<script setup lang="ts">
import { useUser } from "@/stores/useUser"

import {
	useNotification,
	NSpace,
	NTable,
	NPopconfirm,
	NButton,
	NIcon,
} from "naive-ui"

import {
	TrashOutline as DeleteIcon,
} from '@vicons/ionicons5'

const notification = useNotification();

useUser().retrieveUsers();

const deleteUser = (id: string) => {
	useUser().deleteUser(id).then(resp => {
		if (resp.error) {
			notification["error"]({
				content: resp.error,
				duration: 2500,
				keepAliveOnHover: true
			})
			return;
		} 

		notification["success"]({
			content: "Succesfully deleted a user",
			duration: 2500,
			keepAliveOnHover: true
		})
	});
}
</script>

<template>
	<n-space justify="space-between">
		<div class="page-title">Users</div>
	</n-space>
	<n-table :single-line="false">
		<thead>
			<th>ID</th>
			<th>Name</th>
			<th>Created At</th>
			<th>Admin</th>
			<th>Actions</th>
		</thead>
		<tbody>
			<tr v-for="(user, index) in useUser().users" :key="index">
				<td>{{ user.id }}</td>
				<td>{{ user.name }}</td>
				<td>{{ user.created_at }}</td>
				<td>{{ user.admin }}</td>
				<td>
				<n-popconfirm width="310px" placement="top-end" :positive-text="null" negative-text="Confirm" @negative-click="deleteUser(user.id)">
					<template #trigger>
						<n-button type="error" :focusable="false" text>
							<template #icon>
								<n-icon size="20px" color="#d03050">
									<delete-icon />
								</n-icon>
							</template>
							Delete
						</n-button>
					</template>
					Are you sure you want to delete the user?
				</n-popconfirm>
				</td>
			</tr>
		</tbody>
	</n-table>
</template>

<style lang="scss" scoped>
@use "../assets/main.scss";
</style>
