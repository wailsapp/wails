interface WailsLogoProps {
  className?: string;
  size?: number;
}

export default function WailsLogo({ className = '', size = 240 }: WailsLogoProps) {
  return (
    <img
      src="/wails-logo.png"
      alt="Wails"
      width={size}
      height={size}
      className={`object-contain ${className}`}
      style={{
        filter: 'drop-shadow(0 0 60px rgba(239, 68, 68, 0.4))',
      }}
    />
  );
}
