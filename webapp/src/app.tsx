import React from 'react';

import Rhs from 'containers/Rhs';

// Global styles
import 'styles/main.scss'
/**
 * Mattermost plugin allows registering only one component in RHS
 * So, we would be grouping all the different components inside "Rhs" component to generate one final component for registration
 */
const App = (): JSX.Element => <Rhs/>;

export default App;
