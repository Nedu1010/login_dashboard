import { useState } from 'react';
import type { FormEvent } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { authAPI } from '../api/auth';

export default function Login() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      await authAPI.login({ email, password });
      navigate('/dashboard');
    } catch (err: any) {
      setError(err.response?.data?.error || 'Login failed');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className='flex-center'>
      <div className='bg-shapes'>
        <div className='shape shape-1'></div>
        <div className='shape shape-2'></div>
        <div className='shape shape-3'></div>
        <div className='shape shape-4'></div>
      </div>

      <div className='glass-card fade-in' style={{ maxWidth: '420px', width: '100%' }}>
        <div className='text-center mb-3'>
          <div style={{ fontSize: '2.5rem', marginBottom: '0.5rem' }}>🛡️</div>
          <h1 style={{ fontSize: '2rem', marginBottom: '0.5rem' }}>Welcome Back</h1>
          <p className='text-secondary'>Sign in to your account</p>
        </div>

        {error && (
          <div
            className='text-error text-center mb-2'
            style={{
              padding: '0.75rem',
              background: 'rgba(245, 101, 101, 0.1)',
              borderRadius: 'var(--radius-sm)',
              border: '1px solid rgba(245, 101, 101, 0.3)',
            }}
          >
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit}>
          <div className='input-group'>
            <label htmlFor='email'>Email Address</label>
            <input
              id='email'
              type='email'
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder='you@example.com'
              required
            />
          </div>

          <div className='input-group'>
            <label htmlFor='password'>Password</label>
            <input
              id='password'
              type='password'
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder='••••••••'
              required
            />
          </div>

          <div className='flex-between mb-2' style={{ fontSize: '0.875rem' }}>
            <label style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', cursor: 'pointer' }}>
              <input type='checkbox' style={{ width: 'auto' }} />
              <span className='text-secondary'>Remember me</span>
            </label>
            <a href='#'>Forgot password?</a>
          </div>

          <button type='submit' className='btn btn-primary' style={{ width: '100%' }} disabled={loading}>
            {loading ? 'Signing in...' : 'Sign In'}
          </button>
        </form>

        <p className='text-center mt-3 text-secondary'>
          Don't have an account? <Link to='/register'>Sign up</Link>
        </p>
      </div>
    </div>
  );
}
