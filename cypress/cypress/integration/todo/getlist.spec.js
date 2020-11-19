describe('Get Lists from Root', function (){
  it('retrieves all todo lists', function (){
    cy.request({
      method: 'GET',
      failOnStatusCode: false,
      log: true,
      headers: {
        'accept': 'application/json',
      },
      url: 'http://localhost:3001/',
      response: []
    }).then((response) => {
      expect(response.body).to.not.be.null
      assert.equal(response.status, 200)
      assert.equal(response.body.message, "Ok")
      expect(response.body.data).to.not.be.null

      const lists = response.body.data
      lists.forEach((list) => {
        expect(list.id).to.not.be.null
        expect(list.name).to.not.be.null
        expect(list.description).to.not.be.null

        const data = list.list
        data.forEach((todo) => {
          expect(todo.id).to.not.be.null
          expect(todo.description).to.not.be.null
          expect(todo.title).to.not.be.null
        })
      })
    })
  })
})

describe('Get Lists from /list', function (){
  it('retrieves all todo lists', function (){
    cy.request({
      method: 'GET',
      failOnStatusCode: false,
      log: true,
      headers: {
        'accept': 'application/json',
      },
      url: 'http://localhost:3001/list',
      response: []
    }).then((response) => {
      expect(response.body).to.not.be.null
      assert.equal(response.status, 200)
      assert.equal(response.body.message, "Ok")
      expect(response.body.data).to.not.be.null

      const lists = response.body.data
      lists.forEach((list) => {
        expect(list.id).to.not.be.null
        expect(list.name).to.not.be.null
        expect(list.description).to.not.be.null
        
        const data = list.list
        data.forEach((todo) => {
          expect(todo.id).to.not.be.null
          expect(todo.description).to.not.be.null
          expect(todo.title).to.not.be.null
        })
      })
    })
  })
})

describe('Get Lists by ID', function (){
  it('retrieves a todo list given its ID', function (){
    cy.fixture('todo/listids.json').then((ids) => {
      ids.forEach((id)=>{
        cy.request({
          method: 'GET',
          failOnStatusCode: false,
          log: true,
          headers: {
            'accept': 'application/json',
          },
          url: 'http://localhost:3001/list/'+id,
          response: []
        }).then((response) => {
          expect(response.body).to.not.be.null
          assert.equal(response.status, 200)
          assert.equal(response.body.message, "Ok")
          expect(response.body.data).to.not.be.null

          const lists = response.body.data
          lists.forEach((list) => {
            expect(list.id).to.not.be.null
            expect(list.name).to.not.be.null
            expect(list.description).to.not.be.null
            
            const data = list.list
            data.forEach((todo) => {
              expect(todo.id).to.not.be.null
              expect(todo.description).to.not.be.null
              expect(todo.title).to.not.be.null
            })
          })
        })
      })
    })
  })
})
