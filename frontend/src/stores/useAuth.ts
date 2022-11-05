import { ref } from "vue"
import { defineStore } from "pinia"
import { useUser } from "@/stores/useUser"
import axios from 'axios'

export const useAuth = defineStore("auth", () => {
	const key = ref("");
	const user = ref(null);
	const authorized = ref(false);

	const storage_key = localStorage.getItem("key");
	if (storage_key != "") {
		axios
			.get("/api/self", {
				headers: {
					Authorization: 'Bearer ' + storage_key,
				}
			})
			.then((response) => {
				if (response.status == 200) {
					authorized.value = true;
					key.value = storage_key;
					user.value = response.data;
				}
			})
			.catch((error) => {
				clear();
			})
	}

	function clear() {
		authorized.value = false;
		key.value = "";
		user.value = null;
		localStorage.setItem("key", "");
	}

	function authorize(name: string, password: string) {
		return axios
			.post("/api/login", {
				name: name,
				password: password,
			})
			.then((response) => {
				if (response.status == 200) {
					authorized.value = true;
					key.value = response.data.access_token;
					user.value = response.data.user;
					localStorage.setItem("key", key.value);
				}

				return {
					data: response.data,
				}
			})
			.catch((error) => {
				if (error.response.status == 401) {
					return {
						error: "Invalid name or password",
					}
				}

				return {
					error: error.message,
				}
			})
	}

	function register(name: string, password: string) {
		return axios
			.post("/api/register", {
				name: name,
				password: password,
			})
			.then((response) => {
				if (response.status == 200) {
					authorized.value = true;
					console.log(response.data);
					key.value = response.data.access_token;
					user.value = response.data.user;
					localStorage.setItem("key", key.value);
				}

				return {
					data: response.data,
				}
			})
			.catch((error) => {
				if (error.response.status == 409) {
					return {
						error: "User with this name already exists",
					}
				}

				return {
					error: error.message,
				}
			})
	}

	function createAdmin(name: string, password: string) {
		return axios
			.post("/api/users", {
				name: name,
				password: password,
			}, {
				headers: {
					Authorization: 'Bearer ' + key.value,
				}
			})
			.then((response) => {
				useUser().retrieveUsers();
				return {}
			})
			.catch((error) => {
				if (error.response.status == 409) {
					return {
						error: "User with this name already exists",
					}
				}

				return {
					error: error.message,
				}
			})
	}

	function changePassword(password: string, old_password: string) {
		return axios
			.patch("/api/users", {
				old_password: old_password,
				password: password,
			}, {
				headers: {
					Authorization: 'Bearer ' + key.value,
				}
			})
			.then((response) => {
				return {}
			})
			.catch((error) => {
				if (error.response.status == 400) {
					return {
						error: "Invalid old password",
					}
				}

				return {
					error: error.message,
				}
			})
	}

	function deleteAccount() {
		return axios
			.delete("/api/users", {
				headers: {
					Authorization: 'Bearer ' + key.value,
				}
			})
			.then((response) => {
				return {}
			})
			.catch((error) => {
				return {
					error: error.message,
				}
			})
	}

	return { authorized, user, key, authorize, register, clear, deleteAccount, changePassword, createAdmin }
})
