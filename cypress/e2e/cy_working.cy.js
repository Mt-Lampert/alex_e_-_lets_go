const DEVSERVER = 'http://localhost:3000'

describe('Golang Server', () => {
  it('is up and running', () => {
    cy.visit(DEVSERVER)
  })

  it('answers to a ping', () => {
    cy.visit(`${DEVSERVER}/ping`)
    cy.get("#ping").contains("Ping")
  })

  it('has a working navbar', () => {
    cy.visit(`${DEVSERVER}/ping`)
    cy.get('nav').contains("Login")
  })
})
