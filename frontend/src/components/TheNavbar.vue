<script setup lang="ts">
import { ref, computed, h } from "vue"
import { useRouter, RouterLink } from "vue-router"
import { useAuth } from "@/stores/useAuth"
import { 
	useNotification,
	NForm,
	NFormItem,
	NDropdown,
	NInput,
	NSpace,
	NButton,
	NIcon,
	NModal,
} from "naive-ui"
import {
	LogOutOutline as LogoutIcon,
	PersonOutline as ProfileIcon,
	PeopleOutline as UsersIcon,
	DocumentOutline as PlanIcon,
	NutritionOutline as ProductIcon,
	FastFoodOutline as RecipeIcon,
	AppsOutline as HamburgerIcon,
} from '@vicons/ionicons5'

const loginFormRef = ref<Object>(null);
const registerFormRef = ref<Object>(null);
const passwordChangeFormRef = ref<Object>(null);
const createAdminFormRef = ref<Object>(null);

const renderIcon = (icon: Component) => {
	return () => {
		return h(NIcon, null, {
			default: () => h(icon)
		})
	}
}

const options = computed(() => {
	const values = [
		{
			label: 'Products',
			key: 'products',
			icon: renderIcon(ProductIcon),
		},
		{
			label: 'Recipes',
			key: 'recipes',
			icon: renderIcon(RecipeIcon),
		},
		{
			label: 'Plans',
			key: 'plans',
			icon: renderIcon(PlanIcon),
		},
		,
	];
	
	if (useAuth().authorized && useAuth().user.admin) {
		values.push({
			label: 'Users',
			key: 'users',
			icon: renderIcon(UsersIcon),
		})
	}

	return values;
});

const authRules = {
	name: {
		message: "Please input a name",
		max: 63,
		required: true,
	},
	password: {
		message: "The password length must be in a range of 3 to 63 characters",
		min: 4,
		max: 63,
		required: true,
	},
}

const passwordChangeRules = {
	old_password: {
		message: "The old password length must be in a range of 3 to 63 characters",
		min: 4,
		max: 63,
		required: true,
	},
	password: {
		message: "The new password length must be in a range of 3 to 63 characters",
		min: 4,
		max: 63,
		required: true,
	},
}

const loginModalVisible = ref<boolean>(false);
const registerModalVisible = ref<boolean>(false);
const profileModalVisible = ref<boolean>(false);

const notification = useNotification();
const router = useRouter();

const info = ref<Object>({
	name: "",
	password: "",
});

const password_change = ref<Object>({
	old_password: "",
	password: "",
});

const modalError = ref<string>("");

const showLoginModal = () => {
	modalError.value = "";
	info.value = {
		name: "",
		password: "",
	}

	loginModalVisible.value = true;
}

const showRegisterModal = () => {
	modalError.value = "";
	info.value = {
		name: "",
		password: "",
	}

	registerModalVisible.value = true;
}

const showProfileModal = () => {
	modalError.value = "";

	password_change.value = {
		name: "",
		password: "",
	}

	info.value = {
		name: "",
		password: "",
	}

	profileModalVisible.value = true;
}

const authorize = () => {
	loginFormRef.value.validate((errors) => {
		if (errors) {
			return
		}

		useAuth().authorize(info.value.name, info.value.password).then(resp => {
			if (resp.error) {
				modalError.value = resp.error;
				return;
			} 

			notification["success"]({
				content: "Succesfully logged in",
				duration: 2500,
				keepAliveOnHover: true
			})
			loginModalVisible.value = false;
		});
	}).catch((errors) => {
		console.log(errors);
	})
}

const register = () => {
	registerFormRef.value.validate((errors) => {
		if (errors) {
			return
		}

		useAuth().register(info.value.name, info.value.password).then(resp => {
			if (resp.error) {
				modalError.value = resp.error;
				return;
			} 

			notification["success"]({
				content: "Succesfully registered a new user",
				duration: 2500,
				keepAliveOnHover: true
			})
			registerModalVisible.value = false;
		});
	}).catch((errors) => {
		console.log(errors);
	})
}

const changePassword = () => {
	passwordChangeFormRef.value.validate((errors) => {
		if (errors) {
			return
		}

		useAuth().changePassword(password_change.value.password, password_change.value.old_password).then(resp => {
			if (resp.error) {
				modalError.value = resp.error;
				return;
			} 

			notification["success"]({
				content: "Succesfully changed account password",
				duration: 2500,
				keepAliveOnHover: true
			})
			
			password_change.value = {
				old_password: "",
				password: "",
			}
		});
	}).catch((errors) => {
		console.log(errors);
	})
}

const createAdmin = () => {
	createAdminFormRef.value.validate((errors) => {
		if (errors) {
			return
		}

		useAuth().createAdmin(info.value.name, info.value.password).then(resp => {
			if (resp.error) {
				modalError.value = resp.error;
				return;
			} 

			notification["success"]({
				content: "Succesfully created a new admin",
				duration: 2500,
				keepAliveOnHover: true
			})

			info.value = {
				name: "",
				password: "",
			}
		});
	}).catch((errors) => {
		console.log(errors);
	})
}

const deleteAccount = () => {
	useAuth().deleteAccount().then(resp => {
		if (resp.error) {
			modalError.value = resp.error;
			return;
		} 

		notification["success"]({
			content: "Succesfully deleted the account",
			duration: 2500,
			keepAliveOnHover: true
		})

		profileModalVisible.value = false;
		setTimeout(() => {
			useAuth().clear();
		}, 1000) 
	});
}

const logout = () => {
	useAuth().clear()

	notification["success"]({
		content: "Succesfully logged out",
		duration: 2500,
		keepAliveOnHover: true
	})

	if (router.currentRoute.value.name == "users") {
		router.push({
			name: "products"
		})
	}
}

const handleRoute = (key: string) => {
	router.push({
		name: key,
	})
}
</script>

<template>
	<n-space align="center" justify="space-between">
		<div>
			<div class="full-navigation">
				<router-link to="/products" class="router-link">
					<n-button :focusable="false" text>
						<template #icon>
							<n-icon size="20px" color="#36ad6a">
								<product-icon />
							</n-icon>
						</template>
						Products
					</n-button>
				</router-link>
				<router-link to="/recipes" class="router-link">
					<n-button :focusable="false" text>
						<template #icon>
							<n-icon size="20px" color="#36ad6a">
								<recipe-icon />
							</n-icon>
						</template>
						Recipes
					</n-button>
				</router-link>
				<router-link to="/plans" class="router-link">
					<n-button :focusable="false" text>
						<template #icon>
							<n-icon size="20px" color="#36ad6a">
								<plan-icon />
							</n-icon>
						</template>
						Plans
					</n-button>
				</router-link>
				<router-link v-if="useAuth().authorized && useAuth().user.admin" to="/users" class="router-link">
					<n-button :focusable="false" text>
						<template #icon>
							<n-icon size="20px" color="#36ad6a">
								<users-icon />
							</n-icon>
						</template>
						Users
					</n-button>
				</router-link>
			</div>
			<div class="partial-navigation">
				<n-dropdown class="partial-navigation" :options="options" @select="handleRoute">
					<n-icon size="30px" color="#36ad6a">
						<hamburger-icon />
					</n-icon>
				</n-dropdown>
			</div>
		</div>
		<div v-if="!useAuth().authorized" >
			<n-button :focusable="false" @click="showLoginModal()" text>
				<template #icon>
					<n-icon size="20px" color="#36ad6a">
						<logout-icon />
					</n-icon>
				</template>
				Login
			</n-button>
			<n-button :focusable="false" @click="showRegisterModal()" style="margin-left:10px" text>
				<template #icon>
					<n-icon size="20px" color="#36ad6a">
						<logout-icon />
					</n-icon>
				</template>
				Register
			</n-button>
		</div>
		<div v-else>
			<n-button :focusable="false" @click="showProfileModal()" text>
				<template #icon>
					<n-icon size="20px" :color="useAuth().user.admin ? '#d03050' : '#36ad6a'">
						<profile-icon />
					</n-icon>
				</template>
				Profile ({{ useAuth().user.name }})
			</n-button>
			<n-button :focusable="false" @click="logout" style="margin-left:10px" text>
				<template #icon>
					<n-icon size="20px" color="#36ad6a">
						<logout-icon />
					</n-icon>
				</template>
				Logout
			</n-button>
		</div>
	</n-space>
	<div class="divider"></div>
	<n-modal
		v-model:show="loginModalVisible" 
		preset="dialog" 
		:show-icon="false"
	>
		<div class="modal-title">
			Login
		</div>
		<div class="modal-content">
			<div v-if="modalError" style="color: #d03050; margin-bottom: 20px">{{ modalError }}</div>
			<n-form ref="loginFormRef" :model="info" :rules="authRules">
				<n-form-item label="Name" path="name" required>
					<n-input type="text" maxlength="63" v-model:value="info.name" placeholder="Name" show-count />
				</n-form-item>
				<n-form-item label="Password" path="password" required>
					<n-input type="password" v-model:value="info.password" minlength="4" maxlength="63" placeholder="Password" show-password-on="click" />
				</n-form-item>
			</n-form>
		</div>
		<template #action>
			<n-button type="primary" @click="authorize">
				Login
			</n-button>
		</template>
	</n-modal>
	<n-modal
		v-model:show="registerModalVisible" 
		preset="dialog" 
		:show-icon="false"
	>
		<div class="modal-title">
			Register
		</div>
		<div class="modal-content">
			<div v-if="modalError" style="color: #d03050; margin-bottom: 20px">{{ modalError }}</div>
			<n-form ref="registerFormRef" :model="info" :rules="authRules">
				<n-form-item label="Name" path="name" required>
					<n-input type="text" maxlength="63" v-model:value="info.name" placeholder="Name" show-count />
				</n-form-item>
				<n-form-item label="Password" path="password" required>
					<n-input type="password" v-model:value="info.password" minlength="4" maxlength="63" placeholder="Password" show-password-on="click" />
				</n-form-item>
			</n-form>
		</div>
		<template #action>
			<n-button type="primary" @click="register">
				Register
			</n-button>
		</template>
	</n-modal>
	<n-modal
		v-model:show="profileModalVisible" 
		preset="dialog" 
		:show-icon="false"
	>
		<div class="modal-title">
			{{ useAuth().user.name }}
		</div>
		<div style="margin-bottom: 20px">Joined: {{ useAuth().user.created_at }}</div>
		<div class="modal-content">
			<div v-if="modalError" style="color: #d03050; margin-bottom: 20px">{{ modalError }}</div>
			<div class="modal-content-title">Update password</div>
			<n-form ref="passwordChangeFormRef" :model="password_change" :rules="passwordChangeRules">
				<n-form-item label="Old Password" path="old_password" required>
					<n-input type="password" v-model:value="password_change.old_password" minlength="4" maxlength="63" placeholder="Old Password" show-password-on="click" />
				</n-form-item>
				<n-form-item label="New Password" path="password" required>
					<n-input type="password" v-model:value="password_change.password" minlength="4" maxlength="63" placeholder="New Password" show-password-on="click" />
				</n-form-item>
			</n-form>
			<n-space justify="end">
				<n-button type="primary" @click="changePassword()">
					Update
				</n-button>
			</n-space>
			<div v-if="useAuth().user.admin">
				<div class="divider modal-divider"></div>
				<div class="modal-content-title">Create an admin user</div>
				<n-form ref="createAdminFormRef" :model="info" :rules="authRules">
					<n-form-item label="Name" path="name" required>
						<n-input type="text" maxlength="63" v-model:value="info.name" placeholder="Name" show-count />
					</n-form-item>
					<n-form-item label="Password" path="password" required>
						<n-input type="password" v-model:value="info.password" minlength="4" maxlength="63" placeholder="Password" show-password-on="click" />
					</n-form-item>
				</n-form>
				<n-space justify="end">
					<n-button type="primary" @click="createAdmin()">
						Create
					</n-button>
				</n-space>
			</div>
			<div class="divider modal-divider"></div>
			<div class="modal-content-title" style="text-align:center">Delete your account (action non recoverable)</div>
			<n-space justify="center">
				<n-button type="error" @click="deleteAccount()">
					Delete
				</n-button>
			</n-space>
		</div>
	</n-modal>
</template>

<style lang="scss" scoped>
@use "../assets/main.scss";

.partial-navigation {
	display: none;
}

@media screen and (max-width: 700px) {
	.full-navigation {
		display: none;
	}

	.partial-navigation {
		margin-top: 10px;
		margin-left: 10px;
		display: block;
	}
}

.modal-title {
	font-size: 28px;
	color: main.$green;
	font-weight: 600;
	margin-bottom: 10px;
}

.modal-content-title {
	font-size: 17px;
	margin-bottom: 10px;
	color: main.$green;
	font-weight: 600;
}

.router-link {
	text-decoration: none;
	margin: 0 10px;
}

.modal-divider {
	margin: 20px 0;
}
</style>
