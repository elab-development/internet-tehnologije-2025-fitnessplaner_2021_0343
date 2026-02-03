import React from 'react';

interface CardProps {
  children: React.ReactNode;
  title?: string;
  className?: string;
  headerAction?: React.ReactNode;
  footer?: React.ReactNode;
}

const Card: React.FC<CardProps> = ({
  children,
  title,
  className = '',
  headerAction,
  footer,
}) => {
  return (
    <div className={`bg-white rounded-xl shadow-md p-6 ${className}`}>
      {title && (
        <div className="flex justify-between items-center mb-4 pb-4 border-b">
          <h2 className="text-2xl font-bold text-gray-800">{title}</h2>
          {headerAction && <div>{headerAction}</div>}
        </div>
      )}
      <div className="card-content">
        {children}
      </div>
      {footer && (
        <div className="mt-4 pt-4 border-t">
          {footer}
        </div>
      )}
    </div>
  );
};

export default Card;

