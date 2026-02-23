import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { authAPI } from '../api/auth';
import type { User } from '../api/auth';

export default function Dashboard() {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();

  useEffect(() => {
    fetchUser();
  }, []);

  const fetchUser = async () => {
    try {
      const response = await authAPI.getMe();
      setUser(response.data.user);
    } catch (error) {
      navigate('/login');
    } finally {
      setLoading(false);
    }
  };

  const handleLogout = async () => {
    try {
      await authAPI.logout();
      navigate('/login');
    } catch (error) {
      console.error('Logout failed:', error);
    }
  };

  if (loading) {
    return (
      <div className='flex-center'>
        <div style={{ fontSize: '2rem' }}>Loading...</div>
      </div>
    );
  }

  return (
    <div style={{ minHeight: '100vh', padding: '2rem' }}>
      <div className='bg-shapes'>
        <div className='shape shape-1'></div>
        <div className='shape shape-2'></div>
        <div className='shape shape-3'></div>
        <div className='shape shape-4'></div>
      </div>

      <div className='container' style={{ maxWidth: '1000px' }}>
        {/* Header */}
        <div className='glass-card mb-3'>
          <div className='flex-between'>
            <div>
              <h1 style={{ fontSize: '2rem', marginBottom: '0.5rem' }}>Dashboard</h1>
              <p className='text-secondary'>Welcome back, {user?.email}!</p>
            </div>
            <button onClick={handleLogout} className='btn btn-secondary'>
              Logout
            </button>
          </div>
        </div>

        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(300px, 1fr))', gap: '1.5rem' }}>
          {/* Profile Card */}
          <div className='glass-card'>
            <h2
              style={{ fontSize: '1.5rem', marginBottom: '1rem', display: 'flex', alignItems: 'center', gap: '0.5rem' }}
            >
              👤 Profile
            </h2>
            <div style={{ marginBottom: '1rem' }}>
              <p className='text-secondary' style={{ fontSize: '0.875rem', marginBottom: '0.25rem' }}>
                Email
              </p>
              <p style={{ fontSize: '1.125rem' }}>{user?.email}</p>
            </div>
            <div style={{ marginBottom: '1rem' }}>
              <p className='text-secondary' style={{ fontSize: '0.875rem', marginBottom: '0.25rem' }}>
                Status
              </p>
              <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                <span
                  style={{
                    width: '8px',
                    height: '8px',
                    borderRadius: '50%',
                    background: user?.verified ? 'var(--success-green)' : 'var(--error-red)',
                  }}
                ></span>
                <span>{user?.verified ? 'Verified' : 'Not Verified'}</span>
              </div>
            </div>
            <div>
              <p className='text-secondary' style={{ fontSize: '0.875rem', marginBottom: '0.25rem' }}>
                Member Since
              </p>
              <p>{user?.created_at ? new Date(user.created_at).toLocaleDateString() : 'N/A'}</p>
            </div>
          </div>

          {/* Security Card */}
          <div className='glass-card'>
            <h2
              style={{ fontSize: '1.5rem', marginBottom: '1rem', display: 'flex', alignItems: 'center', gap: '0.5rem' }}
            >
              🛡️ Security
            </h2>
            <div style={{ marginBottom: '1rem' }}>
              <p className='text-secondary' style={{ fontSize: '0.875rem', marginBottom: '0.25rem' }}>
                Authentication
              </p>
              <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                <span
                  style={{
                    width: '8px',
                    height: '8px',
                    borderRadius: '50%',
                    background: 'var(--success-green)',
                  }}
                ></span>
                <span>JWT Token Active</span>
              </div>
            </div>
            <div style={{ marginBottom: '1rem' }}>
              <p className='text-secondary' style={{ fontSize: '0.875rem', marginBottom: '0.25rem' }}>
                CSRF Protection
              </p>
              <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                <span
                  style={{
                    width: '8px',
                    height: '8px',
                    borderRadius: '50%',
                    background: 'var(--success-green)',
                  }}
                ></span>
                <span>Enabled</span>
              </div>
            </div>
            <div>
              <p className='text-secondary' style={{ fontSize: '0.875rem', marginBottom: '0.25rem' }}>
                Cookies
              </p>
              <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                <span
                  style={{
                    width: '8px',
                    height: '8px',
                    borderRadius: '50%',
                    background: 'var(--success-green)',
                  }}
                ></span>
                <span>HTTP-Only & Secure</span>
              </div>
            </div>
          </div>
        </div>

        {/* Stats Cards */}
        <div
          style={{
            display: 'grid',
            gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))',
            gap: '1rem',
            marginTop: '1.5rem',
          }}
        >
          <div className='glass-card text-center'>
            <div style={{ fontSize: '2.5rem', marginBottom: '0.5rem' }}>🔐</div>
            <div style={{ fontSize: '2rem', fontWeight: 'bold', marginBottom: '0.25rem' }}>100%</div>
            <div className='text-secondary'>Security Score</div>
          </div>

          <div className='glass-card text-center'>
            <div style={{ fontSize: '2.5rem', marginBottom: '0.5rem' }}>⚡</div>
            <div style={{ fontSize: '2rem', fontWeight: 'bold', marginBottom: '0.25rem' }}>1</div>
            <div className='text-secondary'>Active Session</div>
          </div>

          <div className='glass-card text-center'>
            <div style={{ fontSize: '2.5rem', marginBottom: '0.5rem' }}>✅</div>
            <div style={{ fontSize: '2rem', fontWeight: 'bold', marginBottom: '0.25rem' }}>5m</div>
            <div className='text-secondary'>Token Expiry</div>
          </div>
        </div>
      </div>
    </div>
  );
}
