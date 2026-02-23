import { useState } from "react";
import type { FormEvent } from "react";
import { useNavigate, Link } from "react-router-dom";
import { authAPI } from "../api/auth";

export default function Register() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [error, setError] = useState("");
  const [success, setSuccess] = useState(false);
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const getPasswordStrength = (pwd: string): { text: string; color: string; width: string } => {
    if (pwd.length === 0) return { text: "", color: "", width: "0%" };
    if (pwd.length < 8) return { text: "Weak", color: "#f56565", width: "33%" };
    if (!/(?=.*[a-z])(?=.*[A-Z])(?=.*\d)/.test(pwd)) return { text: "Medium", color: "#ed8936", width: "66%" };
    return { text: "Strong", color: "#48bb78", width: "100%" };
  };

  const passwordStrength = getPasswordStrength(password);

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setError("");
    setSuccess(false);

    if (password !== confirmPassword) {
      setError("Passwords do not match");
      return;
    }

    if (password.length < 8) {
      setError("Password must be at least 8 characters");
      return;
    }

    setLoading(true);

    try {
      await authAPI.register({ email, password });
      setSuccess(true);
      setTimeout(() => navigate("/login"), 2000);
    } catch (err: any) {
      setError(err.response?.data?.error || "Registration failed");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex-center">
      <div className="bg-shapes">
        <div className="shape shape-1"></div>
        <div className="shape shape-2"></div>
        <div className="shape shape-3"></div>
        <div className="shape shape-4"></div>
      </div>

      <div className="glass-card fade-in" style={{ maxWidth: "420px", width: "100%" }}>
        <div className="text-center mb-3">
          <div style={{ fontSize: "2.5rem", marginBottom: "0.5rem" }}>✨</div>
          <h1 style={{ fontSize: "2rem", marginBottom: "0.5rem" }}>Create Account</h1>
          <p className="text-secondary">Join us today</p>
        </div>

        {error && (
          <div
            className="text-error text-center mb-2"
            style={{
              padding: "0.75rem",
              background: "rgba(245, 101, 101, 0.1)",
              borderRadius: "var(--radius-sm)",
              border: "1px solid rgba(245, 101, 101, 0.3)",
            }}
          >
            {error}
          </div>
        )}

        {success && (
          <div
            className="text-success text-center mb-2"
            style={{
              padding: "0.75rem",
              background: "rgba(72, 187, 120, 0.1)",
              borderRadius: "var(--radius-sm)",
              border: "1px solid rgba(72, 187, 120, 0.3)",
            }}
          >
            Registration successful! Redirecting to login...
          </div>
        )}

        <form onSubmit={handleSubmit}>
          <div className="input-group">
            <label htmlFor="email">Email Address</label>
            <input
              id="email"
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="you@example.com"
              required
            />
          </div>

          <div className="input-group">
            <label htmlFor="password">Password</label>
            <input
              id="password"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="••••••••"
              required
            />
            {password && (
              <div style={{ marginTop: "0.5rem" }}>
                <div
                  style={{
                    height: "4px",
                    background: "rgba(255,255,255,0.1)",
                    borderRadius: "2px",
                    overflow: "hidden",
                  }}
                >
                  <div
                    style={{
                      height: "100%",
                      width: passwordStrength.width,
                      background: passwordStrength.color,
                      transition: "all 0.3s ease",
                    }}
                  ></div>
                </div>
                <p style={{ fontSize: "0.75rem", marginTop: "0.25rem", color: passwordStrength.color }}>
                  {passwordStrength.text}
                </p>
              </div>
            )}
          </div>

          <div className="input-group">
            <label htmlFor="confirmPassword">Confirm Password</label>
            <input
              id="confirmPassword"
              type="password"
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
              placeholder="••••••••"
              required
            />
          </div>

          <button type="submit" className="btn btn-primary" style={{ width: "100%" }} disabled={loading}>
            {loading ? "Creating account..." : "Sign Up"}
          </button>
        </form>

        <p className="text-center mt-3 text-secondary">
          Already have an account? <Link to="/login">Sign in</Link>
        </p>
      </div>
    </div>
  );
}
