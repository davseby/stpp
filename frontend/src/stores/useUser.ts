import { ref } from "vue"
import { defineStore } from "pinia"
import { useAuth } from "@/stores/useAuth"
import axios from 'axios'

export const useUser = defineStore("user", () => {
	const users = ref<Array<Object>>([]);

	function retrieveUsers() {
		return axios
			.get("/api/users", {
				headers: {
					Authorization: 'Bearer ' + useAuth().key,
				}
			})
			.then((response) => {
				if (response.status == 200) {
					users.value = response.data;
				}

				return {
					data: response.data,
				}
			})
			.catch((error) => {
				return {
					error: error.message,
				}
			})
	}

	function deleteUser(id: string) {
		return axios
			.delete("/api/users/" + id, {
				headers: {
					Authorization: 'Bearer ' + useAuth().key,
				}
			})
			.then((response) => {
				if (response.status == 204) {
					retrieveUsers();
				}

				return {}
			})
			.catch((error) => {
				return {
					error: error.message,
				}
			})
	}

	return { users, retrieveUsers, deleteUser }
})
