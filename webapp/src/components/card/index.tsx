import React from 'react';

import './styles.scss';

type CardProps = {
    cardHeader: string;
    cardBody: Record<string, string | number>;
    onDelete?: (e: React.SyntheticEvent) => void
}

const Card = ({cardHeader, cardBody, onDelete}: CardProps) => {
    return (
        <div className='details-card'>
            <div className='details-card__header d-flex align-items-center justify-content-between'>
                <h3 className='details-card__header-text'>{cardHeader}</h3>
                {onDelete && (
                    <button
                        onClick={onDelete}
                        className='details-card__delete-btn'
                    >
                        <i className='fa fa-trash'/>
                    </button>
                )}
            </div>
            <ul className='details-card__body'>
                {
                    Object.keys(cardBody).map((header) => (
                        <li
                            key={header}
                            className='body-item'
                        >
                            <p className='body-item__text'>
                                <span className='body-item__header'>{header + ':'}</span>
                                <span>{cardBody[header]}</span>
                            </p>
                        </li>
                    ))
                }
            </ul>
        </div>
    );
};

export default Card;
