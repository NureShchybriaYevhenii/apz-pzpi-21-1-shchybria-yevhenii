import { useNavigate, useSearchParams } from "react-router-dom";
import { Container, Typography } from "@mui/material";
import LoginForm from "../components/LoginForm.jsx";
import RegisterForm from "../components/RegisterForm.jsx";
import { login, register } from "../services/authService.js";
import { useAuthContext } from "../context/AuthContext.jsx";

const AuthPage = () => {
	// eslint-disable-next-line no-unused-vars
	const [_, setAuth] = useAuthContext();
	const navigate = useNavigate();
	const [searchParams] = useSearchParams();
	const authMode = searchParams.get("mode");

	const authenticateUser = async (authData) => {
		let authResponse;

		if (authMode === "login") {
			authResponse = await login(authData);
		} else if (authMode === "register") {
			authResponse = await register(authData);
		} else {
			console.error("Unsupported auth mode");
			return;
		}

		setAuth({
			token: authResponse.tokens.access_token,
			userID: authResponse.user_id,
			role: authResponse.role,
		});

		navigate("/");
	};

	if (authMode === "login") {
		return (
			<Container>
				<LoginForm authenticateUser={authenticateUser} />
			</Container>
		);
	} else if (authMode === "register") {
		return (
			<Container>
				<RegisterForm authenticateUser={authenticateUser} />
			</Container>
		);
	} else {
		return <Typography>Unsupported auth mode</Typography>;
	}
};

export default AuthPage;